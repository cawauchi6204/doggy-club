import 'package:dio/dio.dart';
import 'package:doggyclub/models/dog.dart';

class DogService {
  final Dio _dio;

  DogService({required Dio dio}) : _dio = dio;

  // Dog management
  Future<Dog> createDog(CreateDogRequest request) async {
    try {
      final response = await _dio.post('/api/dogs', data: request.toJson());
      return Dog.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<List<Dog>> getUserDogs() async {
    try {
      final response = await _dio.get('/api/dogs');
      final dogs = (response.data['dogs'] as List)
          .map((dog) => Dog.fromJson(dog))
          .toList();
      return dogs;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<Dog> getDog(String dogId) async {
    try {
      final response = await _dio.get('/api/dogs/$dogId');
      return Dog.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<Dog> updateDog(String dogId, UpdateDogRequest request) async {
    try {
      final response = await _dio.put('/api/dogs/$dogId', data: request.toJson());
      return Dog.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<void> deleteDog(String dogId) async {
    try {
      await _dio.delete('/api/dogs/$dogId');
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<List<Dog>> searchPublicDogs({
    required String query,
    int limit = 20,
    int offset = 0,
  }) async {
    try {
      final response = await _dio.get('/api/dogs/search', queryParameters: {
        'q': query,
        'limit': limit,
        'offset': offset,
      });
      
      final dogs = (response.data['dogs'] as List)
          .map((dog) => Dog.fromJson(dog))
          .toList();
      return dogs;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<List<String>> getPersonalityTraits() async {
    try {
      final response = await _dio.get('/api/dogs/personality-traits');
      return List<String>.from(response.data['traits']);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  DogException _handleError(DioException e) {
    if (e.response?.data != null && e.response?.data['error'] != null) {
      final error = e.response!.data['error'];
      if (error is Map<String, dynamic>) {
        return DogException(error['message'] ?? 'Unknown error');
      }
      return DogException(error.toString());
    }
    return DogException('Network error occurred');
  }
}

class DogException implements Exception {
  final String message;
  
  DogException(this.message);
  
  @override
  String toString() => 'DogException: $message';
}