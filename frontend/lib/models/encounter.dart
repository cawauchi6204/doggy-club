import 'package:freezed_annotation/freezed_annotation.dart';

part 'encounter.freezed.dart';
part 'encounter.g.dart';

@freezed
class Encounter with _$Encounter {
  const factory Encounter({
    required String id,
    required String dog1Id,
    required String dog2Id,
    required List<double> location, // [longitude, latitude]
    required DateTime timestamp,
  }) = _Encounter;

  factory Encounter.fromJson(Map<String, dynamic> json) =>
      _$EncounterFromJson(json);
}

@freezed
class DeviceLocation with _$DeviceLocation {
  const factory DeviceLocation({
    required String id,
    required String dogId,
    required List<double> location, // [longitude, latitude]
    required DateTime timestamp,
    required DateTime updatedAt,
  }) = _DeviceLocation;

  factory DeviceLocation.fromJson(Map<String, dynamic> json) =>
      _$DeviceLocationFromJson(json);
}

// Request DTOs
@freezed
class DetectEncountersRequest with _$DetectEncountersRequest {
  const factory DetectEncountersRequest({
    required String dogId,
    required double radiusMeters,
  }) = _DetectEncountersRequest;

  factory DetectEncountersRequest.fromJson(Map<String, dynamic> json) =>
      _$DetectEncountersRequestFromJson(json);
}

@freezed
class EncounterResponse with _$EncounterResponse {
  const factory EncounterResponse({
    required List<Encounter> encounters,
    required int total,
    required int limit,
    required int offset,
  }) = _EncounterResponse;

  factory EncounterResponse.fromJson(Map<String, dynamic> json) =>
      _$EncounterResponseFromJson(json);
}