import 'package:freezed_annotation/freezed_annotation.dart';

part 'gift.freezed.dart';
part 'gift.g.dart';

@freezed
class Gift with _$Gift {
  const factory Gift({
    required String id,
    required String senderDogId,
    required String receiverDogId,
    required String type,
    required DateTime createdAt,
  }) = _Gift;

  factory Gift.fromJson(Map<String, dynamic> json) => _$GiftFromJson(json);
}

// Request DTOs
@freezed
class SendGiftRequest with _$SendGiftRequest {
  const factory SendGiftRequest({
    required String senderDogId,
    required String receiverDogId,
    required String type,
  }) = _SendGiftRequest;

  factory SendGiftRequest.fromJson(Map<String, dynamic> json) =>
      _$SendGiftRequestFromJson(json);
}

@freezed
class GiftResponse with _$GiftResponse {
  const factory GiftResponse({
    required List<Gift> gifts,
    required int total,
    required int limit,
    required int offset,
  }) = _GiftResponse;

  factory GiftResponse.fromJson(Map<String, dynamic> json) =>
      _$GiftResponseFromJson(json);
}

@freezed
class GiftTypeInfo with _$GiftTypeInfo {
  const factory GiftTypeInfo({
    required String type,
    required String name,
  }) = _GiftTypeInfo;

  factory GiftTypeInfo.fromJson(Map<String, dynamic> json) =>
      _$GiftTypeInfoFromJson(json);
}

@freezed
class GiftRanking with _$GiftRanking {
  const factory GiftRanking({
    required String type,
    required int count,
  }) = _GiftRanking;

  factory GiftRanking.fromJson(Map<String, dynamic> json) =>
      _$GiftRankingFromJson(json);
}

// Constants for simplified gift types
class GiftType {
  static const String bone = 'bone';
  static const String ball = 'ball';
  static const String treat = 'treat';
  static const String toy = 'toy';
  static const String heart = 'heart';
  static const String star = 'star';
  static const String diamond = 'diamond';
  
  static const List<String> all = [
    bone,
    ball,
    treat,
    toy,
    heart,
    star,
    diamond,
  ];
}