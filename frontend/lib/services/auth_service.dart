import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'dart:convert';
import 'package:doggyclub/models/user.dart';
import 'package:doggyclub/models/auth.dart';

class AuthService {
  static const String _accessTokenKey = 'access_token';
  static const String _refreshTokenKey = 'refresh_token';
  static const String _userKey = 'user_data';
  
  final Dio _dio;
  final FlutterSecureStorage _storage;
  
  AuthService({required Dio dio}) 
    : _dio = dio,
      _storage = const FlutterSecureStorage();

  // Authentication endpoints
  Future<AuthResponse> register(RegisterRequest request) async {
    try {
      final response = await _dio.post(
        '/api/auth/register',
        data: request.toJson(),
      );
      
      final authResponse = AuthResponse.fromJson(response.data);
      await _storeTokens(authResponse);
      return authResponse;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<AuthResponse> login(LoginRequest request) async {
    try {
      final response = await _dio.post(
        '/api/auth/login',
        data: request.toJson(),
      );
      
      final authResponse = AuthResponse.fromJson(response.data);
      await _storeTokens(authResponse);
      return authResponse;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<AuthResponse> refreshToken() async {
    final refreshToken = await _storage.read(key: _refreshTokenKey);
    if (refreshToken == null) {
      throw AuthException('No refresh token available');
    }

    try {
      final response = await _dio.post(
        '/api/auth/refresh',
        data: {'refresh_token': refreshToken},
      );
      
      final authResponse = AuthResponse.fromJson(response.data);
      await _storeTokens(authResponse);
      return authResponse;
    } on DioException catch (e) {
      await clearTokens();
      throw _handleError(e);
    }
  }

  Future<void> logout() async {
    final refreshToken = await _storage.read(key: _refreshTokenKey);
    
    try {
      await _dio.post(
        '/api/auth/logout',
        data: {'refresh_token': refreshToken},
      );
    } catch (e) {
      // Continue with logout even if server call fails
    }
    
    await clearTokens();
  }

  Future<void> changePassword({
    required String oldPassword,
    required String newPassword,
  }) async {
    try {
      await _dio.post(
        '/api/auth/change-password',
        data: {
          'old_password': oldPassword,
          'new_password': newPassword,
        },
      );
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<void> forgotPassword({required String email}) async {
    try {
      await _dio.post(
        '/api/auth/forgot-password',
        data: {'email': email},
      );
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<void> resetPassword({
    required String token,
    required String newPassword,
  }) async {
    try {
      await _dio.post(
        '/api/auth/reset-password',
        data: {
          'token': token,
          'new_password': newPassword,
        },
      );
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  // Token management
  Future<String?> getAccessToken() async {
    return await _storage.read(key: _accessTokenKey);
  }

  Future<String?> getToken() async {
    return await getAccessToken();
  }

  Future<String?> getRefreshToken() async {
    return await _storage.read(key: _refreshTokenKey);
  }

  Future<User?> getCurrentUser() async {
    final userData = await _storage.read(key: _userKey);
    if (userData != null) {
      try {
        final json = jsonDecode(userData) as Map<String, dynamic>;
        return User.fromJson(json);
      } catch (e) {
        return null;
      }
    }
    return null;
  }

  Future<bool> isLoggedIn() async {
    final accessToken = await getAccessToken();
    return accessToken != null;
  }

  Future<void> _storeTokens(AuthResponse authResponse) async {
    await _storage.write(key: _accessTokenKey, value: authResponse.token);
    // Store user data as JSON string
    await _storage.write(key: _userKey, value: jsonEncode(authResponse.user.toJson()));
  }

  Future<void> clearTokens() async {
    await _storage.delete(key: _accessTokenKey);
    await _storage.delete(key: _refreshTokenKey);
    await _storage.delete(key: _userKey);
  }

  AuthException _handleError(DioException e) {
    if (e.response?.data != null && e.response?.data['error'] != null) {
      final error = e.response!.data['error'];
      if (error is Map<String, dynamic>) {
        return AuthException(error['message'] ?? 'Unknown error');
      }
      return AuthException(error.toString());
    }
    return AuthException('Network error occurred');
  }
}

// AuthResponse moved to models/auth.dart

class AuthException implements Exception {
  final String message;
  
  AuthException(this.message);
  
  @override
  String toString() => 'AuthException: $message';
}