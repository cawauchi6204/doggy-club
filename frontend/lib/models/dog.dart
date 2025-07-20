import 'package:freezed_annotation/freezed_annotation.dart';

part 'dog.freezed.dart';
part 'dog.g.dart';

@freezed
class Dog with _$Dog {
  const factory Dog({
    required String id,
    required String userId,
    required String name,
    String? breed,
    int? age,
    String? photoUrl,
    String? bio,
    required DateTime createdAt,
  }) = _Dog;

  factory Dog.fromJson(Map<String, dynamic> json) => _$DogFromJson(json);
}

// Vaccination records disabled in simplified schema
// @freezed
// class VaccinationRecord with _$VaccinationRecord {
//   const factory VaccinationRecord({
//     required String id,
//     required String dogId,
//     required String vaccineName,
//     required DateTime dateGiven,
//     DateTime? nextDueDate,
//     required String veterinarian,
//     String? notes,
//     required DateTime createdAt,
//     required DateTime updatedAt,
//   }) = _VaccinationRecord;
//
//   factory VaccinationRecord.fromJson(Map<String, dynamic> json) =>
//       _$VaccinationRecordFromJson(json);
// }

// Request models
@freezed
class CreateDogRequest with _$CreateDogRequest {
  const factory CreateDogRequest({
    required String name,
    String? breed,
    int? age,
    String? photoUrl,
    String? bio,
  }) = _CreateDogRequest;

  factory CreateDogRequest.fromJson(Map<String, dynamic> json) =>
      _$CreateDogRequestFromJson(json);
}

@freezed
class UpdateDogRequest with _$UpdateDogRequest {
  const factory UpdateDogRequest({
    String? name,
    String? breed,
    int? age,
    String? photoUrl,
    String? bio,
  }) = _UpdateDogRequest;

  factory UpdateDogRequest.fromJson(Map<String, dynamic> json) =>
      _$UpdateDogRequestFromJson(json);
}

class PersonalityTraits {
  static const List<String> all = [
    'friendly',
    'playful',
    'energetic',
    'calm',
    'protective',
    'intelligent',
    'loyal',
    'independent',
    'social',
    'gentle',
    'curious',
    'brave',
    'shy',
    'affectionate',
    'alert',
  ];
}