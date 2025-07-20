import 'package:freezed_annotation/freezed_annotation.dart';

part 'user.freezed.dart';
part 'user.g.dart';

@freezed
class User with _$User {
  const factory User({
    required String id,
    required String username,
    required String email,
    @Default('public') String visibility,
    required DateTime createdAt,
  }) = _User;

  factory User.fromJson(Map<String, dynamic> json) => _$UserFromJson(json);
}

// Simplified privacy settings for the new schema
@freezed
class PrivacySettings with _$PrivacySettings {
  const factory PrivacySettings({
    @Default('public') String visibility,
  }) = _PrivacySettings;

  factory PrivacySettings.fromJson(Map<String, dynamic> json) =>
      _$PrivacySettingsFromJson(json);
}

// Request models
@freezed
class UpdateProfileRequest with _$UpdateProfileRequest {
  const factory UpdateProfileRequest({
    String? username,
    String? email,
  }) = _UpdateProfileRequest;

  factory UpdateProfileRequest.fromJson(Map<String, dynamic> json) =>
      _$UpdateProfileRequestFromJson(json);
}

@freezed
class UpdatePrivacySettingsRequest with _$UpdatePrivacySettingsRequest {
  const factory UpdatePrivacySettingsRequest({
    required String visibility,
  }) = _UpdatePrivacySettingsRequest;

  factory UpdatePrivacySettingsRequest.fromJson(Map<String, dynamic> json) =>
      _$UpdatePrivacySettingsRequestFromJson(json);
}