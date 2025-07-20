// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'encounter.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$EncounterImpl _$$EncounterImplFromJson(Map<String, dynamic> json) =>
    _$EncounterImpl(
      id: json['id'] as String,
      dog1Id: json['dog1_id'] as String,
      dog2Id: json['dog2_id'] as String,
      location: (json['location'] as List<dynamic>)
          .map((e) => (e as num).toDouble())
          .toList(),
      timestamp: DateTime.parse(json['timestamp'] as String),
    );

Map<String, dynamic> _$$EncounterImplToJson(_$EncounterImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'dog1_id': instance.dog1Id,
      'dog2_id': instance.dog2Id,
      'location': instance.location,
      'timestamp': instance.timestamp.toIso8601String(),
    };

_$DeviceLocationImpl _$$DeviceLocationImplFromJson(Map<String, dynamic> json) =>
    _$DeviceLocationImpl(
      id: json['id'] as String,
      dogId: json['dog_id'] as String,
      location: (json['location'] as List<dynamic>)
          .map((e) => (e as num).toDouble())
          .toList(),
      timestamp: DateTime.parse(json['timestamp'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );

Map<String, dynamic> _$$DeviceLocationImplToJson(
        _$DeviceLocationImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'dog_id': instance.dogId,
      'location': instance.location,
      'timestamp': instance.timestamp.toIso8601String(),
      'updated_at': instance.updatedAt.toIso8601String(),
    };

_$DetectEncountersRequestImpl _$$DetectEncountersRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$DetectEncountersRequestImpl(
      dogId: json['dog_id'] as String,
      radiusMeters: (json['radius_meters'] as num).toDouble(),
    );

Map<String, dynamic> _$$DetectEncountersRequestImplToJson(
        _$DetectEncountersRequestImpl instance) =>
    <String, dynamic>{
      'dog_id': instance.dogId,
      'radius_meters': instance.radiusMeters,
    };

_$EncounterResponseImpl _$$EncounterResponseImplFromJson(
        Map<String, dynamic> json) =>
    _$EncounterResponseImpl(
      encounters: (json['encounters'] as List<dynamic>)
          .map((e) => Encounter.fromJson(e as Map<String, dynamic>))
          .toList(),
      total: (json['total'] as num).toInt(),
      limit: (json['limit'] as num).toInt(),
      offset: (json['offset'] as num).toInt(),
    );

Map<String, dynamic> _$$EncounterResponseImplToJson(
        _$EncounterResponseImpl instance) =>
    <String, dynamic>{
      'encounters': instance.encounters.map((e) => e.toJson()).toList(),
      'total': instance.total,
      'limit': instance.limit,
      'offset': instance.offset,
    };
