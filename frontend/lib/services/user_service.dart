import 'package:dio/dio.dart';
import 'package:doggyclub/models/user.dart';

class UserService {
  final Dio _dio;

  UserService({required Dio dio}) : _dio = dio;

  // Profile management
  Future<User> getProfile() async {
    try {
      final response = await _dio.get('/api/users/profile');
      return User.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<User> updateProfile({
    String? nickname,
    String? profileImage,
  }) async {
    try {
      final data = <String, dynamic>{};
      if (nickname != null) data['nickname'] = nickname;
      if (profileImage != null) data['profile_image'] = profileImage;

      final response = await _dio.put('/api/users/profile', data: data);
      return User.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<void> updatePrivacySettings({
    bool? shareLocation,
    bool? shareProfile,
    bool? allowMessages,
    List<String>? blockedUserIds,
    List<String>? visibleFields,
  }) async {
    try {
      final data = <String, dynamic>{};
      if (shareLocation != null) data['share_location'] = shareLocation;
      if (shareProfile != null) data['share_profile'] = shareProfile;
      if (allowMessages != null) data['allow_messages'] = allowMessages;
      if (blockedUserIds != null) data['blocked_user_ids'] = blockedUserIds;
      if (visibleFields != null) data['visible_fields'] = visibleFields;

      await _dio.put('/api/users/privacy', data: data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<void> updateNotificationPreferences({
    bool? encounters,
    bool? gifts,
    bool? messages,
    bool? likes,
    bool? comments,
    bool? follows,
    bool? emailDigest,
    bool? pushEnabled,
  }) async {
    try {
      final data = <String, dynamic>{};
      if (encounters != null) data['encounters'] = encounters;
      if (gifts != null) data['gifts'] = gifts;
      if (messages != null) data['messages'] = messages;
      if (likes != null) data['likes'] = likes;
      if (comments != null) data['comments'] = comments;
      if (follows != null) data['follows'] = follows;
      if (emailDigest != null) data['email_digest'] = emailDigest;
      if (pushEnabled != null) data['push_enabled'] = pushEnabled;

      await _dio.put('/api/users/notifications', data: data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<Map<String, dynamic>> getUserCurrency() async {
    try {
      final response = await _dio.get('/api/users/currency');
      return response.data as Map<String, dynamic>;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<void> deleteAccount() async {
    try {
      await _dio.delete('/api/users/profile');
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<Map<String, dynamic>> searchUsers({
    required String query,
    int limit = 20,
    int offset = 0,
  }) async {
    try {
      final response = await _dio.get(
        '/api/users/search',
        queryParameters: {
          'q': query,
          'limit': limit,
          'offset': offset,
        },
      );
      return response.data as Map<String, dynamic>;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  UserException _handleError(DioException e) {
    if (e.response?.data != null && e.response?.data['error'] != null) {
      final error = e.response!.data['error'];
      if (error is Map<String, dynamic>) {
        return UserException(error['message'] ?? 'Unknown error');
      }
      return UserException(error.toString());
    }
    return UserException('Network error occurred');
  }
}

class UserException implements Exception {
  final String message;
  
  UserException(this.message);
  
  @override
  String toString() => 'UserException: $message';
}