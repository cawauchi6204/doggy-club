import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:doggyclub/models/dog.dart';
import 'package:doggyclub/services/dog_service.dart';
import 'package:doggyclub/providers/auth_provider.dart';

// Dog state
class DogsState {
  final List<Dog> dogs;
  final bool isLoading;
  final String? error;

  const DogsState({
    this.dogs = const [],
    this.isLoading = false,
    this.error,
  });

  DogsState copyWith({
    List<Dog>? dogs,
    bool? isLoading,
    String? error,
  }) {
    return DogsState(
      dogs: dogs ?? this.dogs,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

// Dogs notifier
class DogsNotifier extends StateNotifier<DogsState> {
  final DogService _dogService;

  DogsNotifier(this._dogService) : super(const DogsState());

  Future<void> loadUserDogs() async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final dogs = await _dogService.getUserDogs();
      state = state.copyWith(dogs: dogs, isLoading: false);
    } catch (e) {
      state = state.copyWith(
        error: e.toString(),
        isLoading: false,
      );
    }
  }

  Future<void> createDog(CreateDogRequest request) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final dog = await _dogService.createDog(request);

      state = state.copyWith(
        dogs: [...state.dogs, dog],
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        error: e.toString(),
        isLoading: false,
      );
      rethrow;
    }
  }

  Future<void> updateDog(String dogId, UpdateDogRequest request) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final dog = await _dogService.updateDog(dogId, request);

      final updatedDogs = state.dogs.map((d) {
        return d.id == dogId ? dog : d;
      }).toList();

      state = state.copyWith(
        dogs: updatedDogs,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        error: e.toString(),
        isLoading: false,
      );
      rethrow;
    }
  }

  Future<void> deleteDog(String dogId) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      await _dogService.deleteDog(dogId);

      final updatedDogs = state.dogs.where((d) => d.id != dogId).toList();

      state = state.copyWith(
        dogs: updatedDogs,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        error: e.toString(),
        isLoading: false,
      );
      rethrow;
    }
  }

  Future<List<Dog>> searchPublicDogs({
    required String query,
    int limit = 20,
    int offset = 0,
  }) async {
    try {
      return await _dogService.searchPublicDogs(
        query: query,
        limit: limit,
        offset: offset,
      );
    } catch (e) {
      state = state.copyWith(error: e.toString());
      rethrow;
    }
  }

  void clearError() {
    state = state.copyWith(error: null);
  }
}

// Providers
final dogServiceProvider = Provider<DogService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return apiClient.dogService;
});

final dogsProvider = StateNotifierProvider<DogsNotifier, DogsState>((ref) {
  final dogService = ref.watch(dogServiceProvider);
  return DogsNotifier(dogService);
});

// Convenience providers
final userDogsProvider = Provider<List<Dog>>((ref) {
  return ref.watch(dogsProvider).dogs;
});

final isDogsLoadingProvider = Provider<bool>((ref) {
  return ref.watch(dogsProvider).isLoading;
});

// Personality traits provider
final personalityTraitsProvider = FutureProvider<List<String>>((ref) async {
  final dogService = ref.watch(dogServiceProvider);
  return await dogService.getPersonalityTraits();
});