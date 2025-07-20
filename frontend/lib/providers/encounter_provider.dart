import 'dart:async';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:geolocator/geolocator.dart';
import 'package:doggyclub/models/encounter.dart';
import 'package:doggyclub/services/encounter_service.dart';
import 'package:doggyclub/services/location_service.dart';
import 'package:doggyclub/providers/auth_provider.dart';
import 'package:doggyclub/providers/dog_provider.dart';

// Encounter state
class EncounterState {
  final List<Encounter> encounters;
  final bool isLocationTracking;
  final bool isLoading;
  final String? error;
  final Position? currentPosition;

  const EncounterState({
    this.encounters = const [],
    this.isLocationTracking = false,
    this.isLoading = false,
    this.error,
    this.currentPosition,
  });

  EncounterState copyWith({
    List<Encounter>? encounters,
    bool? isLocationTracking,
    bool? isLoading,
    String? error,
    Position? currentPosition,
  }) {
    return EncounterState(
      encounters: encounters ?? this.encounters,
      isLocationTracking: isLocationTracking ?? this.isLocationTracking,
      isLoading: isLoading ?? this.isLoading,
      error: error,
      currentPosition: currentPosition ?? this.currentPosition,
    );
  }
}

// Encounter notifier
class EncounterNotifier extends StateNotifier<EncounterState> {
  final EncounterService _encounterService;
  final LocationService _locationService;
  final Ref _ref;
  Timer? _backgroundDetectionTimer;

  EncounterNotifier(this._encounterService, this._locationService, this._ref) 
      : super(const EncounterState()) {
    _initLocationService();
  }

  void _initLocationService() {
    _locationService.onLocationChanged = _onLocationChanged;
  }

  Future<void> startLocationTracking() async {
    state = state.copyWith(isLoading: true, error: null);
    
    try {
      final success = await _locationService.startLocationTracking();
      if (success) {
        state = state.copyWith(
          isLocationTracking: true,
          isLoading: false,
        );
        
        // Start background detection timer (every 30 seconds)
        _startBackgroundDetection();
      } else {
        state = state.copyWith(
          error: 'Failed to start location tracking. Please check permissions.',
          isLoading: false,
        );
      }
    } catch (e) {
      state = state.copyWith(
        error: e.toString(),
        isLoading: false,
      );
    }
  }

  void stopLocationTracking() {
    _locationService.stopLocationTracking();
    _stopBackgroundDetection();
    state = state.copyWith(
      isLocationTracking: false,
    );
  }

  Future<void> detectEncounters({String? specificDogId}) async {
    final currentPosition = _locationService.getLastKnownPosition();
    if (currentPosition == null) {
      return;
    }

    // Get user's first dog if no specific dog is provided
    String? dogId = specificDogId;
    if (dogId == null) {
      final dogs = _ref.read(userDogsProvider);
      if (dogs.isEmpty) {
        return;
      }
      dogId = dogs.first.id;
    }

    try {
      final request = DetectEncountersRequest(
        dogId: dogId,
        radiusMeters: 100.0, // Default radius
      );

      final response = await _encounterService.detectEncounters(request);
      
      // For simplified schema, we just get encounter history
      await loadEncounterHistory();
    } catch (e) {
      // Silent fail for background detection
      print('Encounter detection failed: $e');
    }
  }

  Future<void> loadEncounterHistory({String? dogId}) async {
    state = state.copyWith(isLoading: true, error: null);
    
    try {
      // Get user's first dog if no specific dog is provided
      String? selectedDogId = dogId;
      if (selectedDogId == null) {
        final dogs = _ref.read(userDogsProvider);
        if (dogs.isEmpty) {
          state = state.copyWith(
            error: 'No dogs found',
            isLoading: false,
          );
          return;
        }
        selectedDogId = dogs.first.id;
      }

      final response = await _encounterService.getEncounterHistory(
        dogId: selectedDogId,
      );
      
      state = state.copyWith(
        encounters: response.encounters,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        error: e.toString(),
        isLoading: false,
      );
    }
  }

  void _onLocationChanged(Position position) {
    state = state.copyWith(currentPosition: position);
  }

  void _startBackgroundDetection() {
    _backgroundDetectionTimer = Timer.periodic(
      const Duration(seconds: 30),
      (timer) {
        if (state.isLocationTracking) {
          detectEncounters();
        }
      },
    );
  }

  void _stopBackgroundDetection() {
    _backgroundDetectionTimer?.cancel();
    _backgroundDetectionTimer = null;
  }

  void clearError() {
    state = state.copyWith(error: null);
  }

  @override
  void dispose() {
    _stopBackgroundDetection();
    _locationService.dispose();
    super.dispose();
  }
}

// Providers
final encounterServiceProvider = Provider<EncounterService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return apiClient.encounterService;
});

final locationServiceProvider = Provider<LocationService>((ref) {
  return LocationService();
});

final encounterProvider = StateNotifierProvider<EncounterNotifier, EncounterState>((ref) {
  final encounterService = ref.watch(encounterServiceProvider);
  final locationService = ref.watch(locationServiceProvider);
  return EncounterNotifier(encounterService, locationService, ref);
});

// Convenience providers
final encounterHistoryProvider = Provider<List<Encounter>>((ref) {
  return ref.watch(encounterProvider).encounters;
});

final isLocationTrackingProvider = Provider<bool>((ref) {
  return ref.watch(encounterProvider).isLocationTracking;
});

final currentPositionProvider = Provider<Position?>((ref) {
  return ref.watch(encounterProvider).currentPosition;
});