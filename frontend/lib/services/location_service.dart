import 'dart:async';
import 'package:geolocator/geolocator.dart';

class LocationService {
  StreamSubscription<Position>? _positionSubscription;
  Position? _lastKnownPosition;
  
  // Location change callback
  Function(Position)? onLocationChanged;
  
  // Encounter detection callback
  Function(Position)? onEncounterDetectionTriggered;
  
  // Last encounter detection time to prevent spam
  DateTime? _lastEncounterDetection;
  static const Duration _encounterCooldown = Duration(minutes: 2);

  // Start location tracking
  Future<bool> startLocationTracking() async {
    // Request location permissions
    final hasPermission = await _requestLocationPermission();
    if (!hasPermission) {
      return false;
    }

    // Check if location services are enabled
    final isEnabled = await Geolocator.isLocationServiceEnabled();
    if (!isEnabled) {
      return false;
    }

    // Start listening to position changes
    const locationSettings = LocationSettings(
      accuracy: LocationAccuracy.high,
      distanceFilter: 10, // Update every 10 meters
    );

    _positionSubscription = Geolocator.getPositionStream(
      locationSettings: locationSettings,
    ).listen(
      _onPositionChanged,
      onError: (error) {
        print('Location error: $error');
      },
    );

    return true;
  }

  // Stop location tracking
  void stopLocationTracking() {
    _positionSubscription?.cancel();
    _positionSubscription = null;
  }

  // Get current position
  Future<Position?> getCurrentPosition() async {
    try {
      final hasPermission = await _requestLocationPermission();
      if (!hasPermission) {
        return null;
      }

      final position = await Geolocator.getCurrentPosition(
        desiredAccuracy: LocationAccuracy.high,
      );
      
      _lastKnownPosition = position;
      return position;
    } catch (e) {
      print('Error getting current position: $e');
      return null;
    }
  }

  // Get last known position
  Position? getLastKnownPosition() {
    return _lastKnownPosition;
  }

  // Calculate distance between two positions
  double calculateDistance(
    double lat1, double lon1,
    double lat2, double lon2,
  ) {
    return Geolocator.distanceBetween(lat1, lon1, lat2, lon2);
  }

  // Check if encounter detection should be triggered
  bool shouldTriggerEncounterDetection() {
    if (_lastEncounterDetection == null) {
      return true;
    }
    
    final now = DateTime.now();
    return now.difference(_lastEncounterDetection!) > _encounterCooldown;
  }

  // Mark encounter detection as triggered
  void markEncounterDetectionTriggered() {
    _lastEncounterDetection = DateTime.now();
  }

  // Internal position change handler
  void _onPositionChanged(Position position) {
    _lastKnownPosition = position;
    
    // Notify location change callback
    onLocationChanged?.call(position);
    
    // Check if we should trigger encounter detection
    if (shouldTriggerEncounterDetection()) {
      onEncounterDetectionTriggered?.call(position);
      markEncounterDetectionTriggered();
    }
  }

  // Request location permissions
  Future<bool> _requestLocationPermission() async {
    // Check current permission status
    LocationPermission permission = await Geolocator.checkPermission();
    
    if (permission == LocationPermission.denied) {
      permission = await Geolocator.requestPermission();
    }
    
    if (permission == LocationPermission.deniedForever) {
      // Permissions are permanently denied, open settings
      await Geolocator.openAppSettings();
      return false;
    }
    
    return permission == LocationPermission.whileInUse || 
           permission == LocationPermission.always;
  }

  // Check location permission status
  Future<LocationPermissionStatus> checkLocationPermissionStatus() async {
    final permission = await Geolocator.checkPermission();
    
    switch (permission) {
      case LocationPermission.always:
        return LocationPermissionStatus.always;
      case LocationPermission.whileInUse:
        return LocationPermissionStatus.whileInUse;
      case LocationPermission.denied:
        return LocationPermissionStatus.denied;
      case LocationPermission.deniedForever:
        return LocationPermissionStatus.deniedForever;
      case LocationPermission.unableToDetermine:
        return LocationPermissionStatus.unknown;
    }
  }

  // Clean up resources
  void dispose() {
    stopLocationTracking();
  }
}

enum LocationPermissionStatus {
  always,
  whileInUse,
  denied,
  deniedForever,
  unknown,
}