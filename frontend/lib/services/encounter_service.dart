import 'package:dio/dio.dart';
import 'package:doggyclub/models/encounter.dart';

class EncounterService {
  final Dio _dio;

  EncounterService({required Dio dio}) : _dio = dio;

  // Encounter detection
  Future<Map<String, dynamic>> detectEncounters(DetectEncountersRequest request) async {
    try {
      final response = await _dio.post('/api/encounters/detect', data: request.toJson());
      return response.data;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  // Encounter history
  Future<EncounterResponse> getEncounterHistory({
    required String dogId,
    int limit = 20,
    int offset = 0,
  }) async {
    try {
      final response = await _dio.get('/api/encounters/history', queryParameters: {
        'dogId': dogId,
        'limit': limit,
        'offset': offset,
      });
      
      return EncounterResponse.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  // Get encounter details
  Future<Encounter> getEncounterDetails(String encounterId) async {
    try {
      final response = await _dio.get('/api/encounters/$encounterId/details');
      return Encounter.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  EncounterException _handleError(DioException e) {
    if (e.response?.data != null && e.response?.data['error'] != null) {
      final error = e.response!.data['error'];
      if (error is Map<String, dynamic>) {
        return EncounterException(error['message'] ?? 'Unknown error');
      }
      return EncounterException(error.toString());
    }
    return EncounterException('Network error occurred');
  }
}

class EncounterException implements Exception {
  final String message;
  
  EncounterException(this.message);
  
  @override
  String toString() => 'EncounterException: $message';
}