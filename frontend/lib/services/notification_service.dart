import 'package:dio/dio.dart';
import 'package:doggyclub/models/notification.dart';

class NotificationService {
  final Dio dio;

  NotificationService({required this.dio});

  // Device management
  Future<Map<String, dynamic>> registerDevice(RegisterDeviceRequest request) async {
    final response = await dio.post(
      '/api/notifications/devices',
      data: request.toJson(),
    );
    return response.data;
  }

  Future<void> unregisterDevice(String deviceToken) async {
    await dio.delete('/api/notifications/devices/$deviceToken');
  }

  // Notifications
  Future<NotificationResponse> getNotifications({
    int limit = 20,
    int offset = 0,
  }) async {
    final response = await dio.get(
      '/api/notifications',
      queryParameters: {
        'limit': limit,
        'offset': offset,
      },
    );
    return NotificationResponse.fromJson(response.data);
  }

  // These methods are disabled in simplified schema but kept for compatibility
  Future<void> markNotificationAsRead(String notificationId) async {
    // Returns 501 Not Implemented in simplified schema
    await dio.put('/api/notifications/$notificationId/read');
  }

  Future<void> markAllNotificationsAsRead() async {
    // Returns 501 Not Implemented in simplified schema
    await dio.put('/api/notifications/read-all');
  }

  Future<int> getUnreadCount() async {
    // Returns 501 Not Implemented in simplified schema
    final response = await dio.get('/api/notifications/unread-count');
    return response.data['unread_count'] as int;
  }
}