import 'package:freezed_annotation/freezed_annotation.dart';

part 'notification.freezed.dart';
part 'notification.g.dart';

@freezed
class AppNotification with _$AppNotification {
  const factory AppNotification({
    required String id,
    required String userId,
    required String title,
    required String message,
    required String type,
    required DateTime createdAt,
  }) = _AppNotification;

  factory AppNotification.fromJson(Map<String, dynamic> json) =>
      _$AppNotificationFromJson(json);
}

@freezed
class NotificationResponse with _$NotificationResponse {
  const factory NotificationResponse({
    required List<AppNotification> notifications,
    required int total,
    required int limit,
    required int offset,
  }) = _NotificationResponse;

  factory NotificationResponse.fromJson(Map<String, dynamic> json) =>
      _$NotificationResponseFromJson(json);
}

// Notification types
class NotificationType {
  static const String encounter = 'encounter';
  static const String gift = 'gift';
  static const String like = 'like';
  static const String comment = 'comment';
  static const String follow = 'follow';
  static const String welcome = 'welcome';
  static const String marketing = 'marketing';
  static const String system = 'system';
}

// Platform types
class PlatformType {
  static const String ios = 'ios';
  static const String android = 'android';
}

// Request DTOs
@freezed
class RegisterDeviceRequest with _$RegisterDeviceRequest {
  const factory RegisterDeviceRequest({
    required String deviceToken,
    required String platform,
    String? deviceInfo,
  }) = _RegisterDeviceRequest;

  factory RegisterDeviceRequest.fromJson(Map<String, dynamic> json) =>
      _$RegisterDeviceRequestFromJson(json);
}