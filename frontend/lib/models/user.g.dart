// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'user.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$UserImpl _$$UserImplFromJson(Map<String, dynamic> json) => _$UserImpl(
      id: json['id'] as String,
      username: json['username'] as String,
      email: json['email'] as String,
      visibility: json['visibility'] as String? ?? 'public',
      createdAt: DateTime.parse(json['created_at'] as String),
    );

Map<String, dynamic> _$$UserImplToJson(_$UserImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'username': instance.username,
      'email': instance.email,
      'visibility': instance.visibility,
      'created_at': instance.createdAt.toIso8601String(),
    };

_$PrivacySettingsImpl _$$PrivacySettingsImplFromJson(
        Map<String, dynamic> json) =>
    _$PrivacySettingsImpl(
      visibility: json['visibility'] as String? ?? 'public',
    );

Map<String, dynamic> _$$PrivacySettingsImplToJson(
        _$PrivacySettingsImpl instance) =>
    <String, dynamic>{
      'visibility': instance.visibility,
    };

_$UpdateProfileRequestImpl _$$UpdateProfileRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$UpdateProfileRequestImpl(
      username: json['username'] as String?,
      email: json['email'] as String?,
    );

Map<String, dynamic> _$$UpdateProfileRequestImplToJson(
        _$UpdateProfileRequestImpl instance) =>
    <String, dynamic>{
      'username': instance.username,
      'email': instance.email,
    };

_$UpdatePrivacySettingsRequestImpl _$$UpdatePrivacySettingsRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$UpdatePrivacySettingsRequestImpl(
      visibility: json['visibility'] as String,
    );

Map<String, dynamic> _$$UpdatePrivacySettingsRequestImplToJson(
        _$UpdatePrivacySettingsRequestImpl instance) =>
    <String, dynamic>{
      'visibility': instance.visibility,
    };
