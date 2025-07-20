// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'gift.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$GiftImpl _$$GiftImplFromJson(Map<String, dynamic> json) => _$GiftImpl(
      id: json['id'] as String,
      senderDogId: json['sender_dog_id'] as String,
      receiverDogId: json['receiver_dog_id'] as String,
      type: json['type'] as String,
      createdAt: DateTime.parse(json['created_at'] as String),
    );

Map<String, dynamic> _$$GiftImplToJson(_$GiftImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'sender_dog_id': instance.senderDogId,
      'receiver_dog_id': instance.receiverDogId,
      'type': instance.type,
      'created_at': instance.createdAt.toIso8601String(),
    };

_$SendGiftRequestImpl _$$SendGiftRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$SendGiftRequestImpl(
      senderDogId: json['sender_dog_id'] as String,
      receiverDogId: json['receiver_dog_id'] as String,
      type: json['type'] as String,
    );

Map<String, dynamic> _$$SendGiftRequestImplToJson(
        _$SendGiftRequestImpl instance) =>
    <String, dynamic>{
      'sender_dog_id': instance.senderDogId,
      'receiver_dog_id': instance.receiverDogId,
      'type': instance.type,
    };

_$GiftResponseImpl _$$GiftResponseImplFromJson(Map<String, dynamic> json) =>
    _$GiftResponseImpl(
      gifts: (json['gifts'] as List<dynamic>)
          .map((e) => Gift.fromJson(e as Map<String, dynamic>))
          .toList(),
      total: (json['total'] as num).toInt(),
      limit: (json['limit'] as num).toInt(),
      offset: (json['offset'] as num).toInt(),
    );

Map<String, dynamic> _$$GiftResponseImplToJson(_$GiftResponseImpl instance) =>
    <String, dynamic>{
      'gifts': instance.gifts.map((e) => e.toJson()).toList(),
      'total': instance.total,
      'limit': instance.limit,
      'offset': instance.offset,
    };

_$GiftTypeInfoImpl _$$GiftTypeInfoImplFromJson(Map<String, dynamic> json) =>
    _$GiftTypeInfoImpl(
      type: json['type'] as String,
      name: json['name'] as String,
    );

Map<String, dynamic> _$$GiftTypeInfoImplToJson(_$GiftTypeInfoImpl instance) =>
    <String, dynamic>{
      'type': instance.type,
      'name': instance.name,
    };

_$GiftRankingImpl _$$GiftRankingImplFromJson(Map<String, dynamic> json) =>
    _$GiftRankingImpl(
      type: json['type'] as String,
      count: (json['count'] as num).toInt(),
    );

Map<String, dynamic> _$$GiftRankingImplToJson(_$GiftRankingImpl instance) =>
    <String, dynamic>{
      'type': instance.type,
      'count': instance.count,
    };
