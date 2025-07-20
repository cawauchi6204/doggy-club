// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'dog.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$DogImpl _$$DogImplFromJson(Map<String, dynamic> json) => _$DogImpl(
      id: json['id'] as String,
      userId: json['user_id'] as String,
      name: json['name'] as String,
      breed: json['breed'] as String?,
      age: (json['age'] as num?)?.toInt(),
      photoUrl: json['photo_url'] as String?,
      bio: json['bio'] as String?,
      createdAt: DateTime.parse(json['created_at'] as String),
    );

Map<String, dynamic> _$$DogImplToJson(_$DogImpl instance) => <String, dynamic>{
      'id': instance.id,
      'user_id': instance.userId,
      'name': instance.name,
      'breed': instance.breed,
      'age': instance.age,
      'photo_url': instance.photoUrl,
      'bio': instance.bio,
      'created_at': instance.createdAt.toIso8601String(),
    };

_$CreateDogRequestImpl _$$CreateDogRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$CreateDogRequestImpl(
      name: json['name'] as String,
      breed: json['breed'] as String?,
      age: (json['age'] as num?)?.toInt(),
      photoUrl: json['photo_url'] as String?,
      bio: json['bio'] as String?,
    );

Map<String, dynamic> _$$CreateDogRequestImplToJson(
        _$CreateDogRequestImpl instance) =>
    <String, dynamic>{
      'name': instance.name,
      'breed': instance.breed,
      'age': instance.age,
      'photo_url': instance.photoUrl,
      'bio': instance.bio,
    };

_$UpdateDogRequestImpl _$$UpdateDogRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$UpdateDogRequestImpl(
      name: json['name'] as String?,
      breed: json['breed'] as String?,
      age: (json['age'] as num?)?.toInt(),
      photoUrl: json['photo_url'] as String?,
      bio: json['bio'] as String?,
    );

Map<String, dynamic> _$$UpdateDogRequestImplToJson(
        _$UpdateDogRequestImpl instance) =>
    <String, dynamic>{
      'name': instance.name,
      'breed': instance.breed,
      'age': instance.age,
      'photo_url': instance.photoUrl,
      'bio': instance.bio,
    };
