// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'dog.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models');

Dog _$DogFromJson(Map<String, dynamic> json) {
  return _Dog.fromJson(json);
}

/// @nodoc
mixin _$Dog {
  String get id => throw _privateConstructorUsedError;
  String get userId => throw _privateConstructorUsedError;
  String get name => throw _privateConstructorUsedError;
  String? get breed => throw _privateConstructorUsedError;
  int? get age => throw _privateConstructorUsedError;
  String? get photoUrl => throw _privateConstructorUsedError;
  String? get bio => throw _privateConstructorUsedError;
  DateTime get createdAt => throw _privateConstructorUsedError;

  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  $DogCopyWith<Dog> get copyWith => throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $DogCopyWith<$Res> {
  factory $DogCopyWith(Dog value, $Res Function(Dog) then) =
      _$DogCopyWithImpl<$Res, Dog>;
  @useResult
  $Res call(
      {String id,
      String userId,
      String name,
      String? breed,
      int? age,
      String? photoUrl,
      String? bio,
      DateTime createdAt});
}

/// @nodoc
class _$DogCopyWithImpl<$Res, $Val extends Dog> implements $DogCopyWith<$Res> {
  _$DogCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? userId = null,
    Object? name = null,
    Object? breed = freezed,
    Object? age = freezed,
    Object? photoUrl = freezed,
    Object? bio = freezed,
    Object? createdAt = null,
  }) {
    return _then(_value.copyWith(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as String,
      userId: null == userId
          ? _value.userId
          : userId // ignore: cast_nullable_to_non_nullable
              as String,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      breed: freezed == breed
          ? _value.breed
          : breed // ignore: cast_nullable_to_non_nullable
              as String?,
      age: freezed == age
          ? _value.age
          : age // ignore: cast_nullable_to_non_nullable
              as int?,
      photoUrl: freezed == photoUrl
          ? _value.photoUrl
          : photoUrl // ignore: cast_nullable_to_non_nullable
              as String?,
      bio: freezed == bio
          ? _value.bio
          : bio // ignore: cast_nullable_to_non_nullable
              as String?,
      createdAt: null == createdAt
          ? _value.createdAt
          : createdAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$DogImplCopyWith<$Res> implements $DogCopyWith<$Res> {
  factory _$$DogImplCopyWith(_$DogImpl value, $Res Function(_$DogImpl) then) =
      __$$DogImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String id,
      String userId,
      String name,
      String? breed,
      int? age,
      String? photoUrl,
      String? bio,
      DateTime createdAt});
}

/// @nodoc
class __$$DogImplCopyWithImpl<$Res> extends _$DogCopyWithImpl<$Res, _$DogImpl>
    implements _$$DogImplCopyWith<$Res> {
  __$$DogImplCopyWithImpl(_$DogImpl _value, $Res Function(_$DogImpl) _then)
      : super(_value, _then);

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? userId = null,
    Object? name = null,
    Object? breed = freezed,
    Object? age = freezed,
    Object? photoUrl = freezed,
    Object? bio = freezed,
    Object? createdAt = null,
  }) {
    return _then(_$DogImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as String,
      userId: null == userId
          ? _value.userId
          : userId // ignore: cast_nullable_to_non_nullable
              as String,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      breed: freezed == breed
          ? _value.breed
          : breed // ignore: cast_nullable_to_non_nullable
              as String?,
      age: freezed == age
          ? _value.age
          : age // ignore: cast_nullable_to_non_nullable
              as int?,
      photoUrl: freezed == photoUrl
          ? _value.photoUrl
          : photoUrl // ignore: cast_nullable_to_non_nullable
              as String?,
      bio: freezed == bio
          ? _value.bio
          : bio // ignore: cast_nullable_to_non_nullable
              as String?,
      createdAt: null == createdAt
          ? _value.createdAt
          : createdAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$DogImpl implements _Dog {
  const _$DogImpl(
      {required this.id,
      required this.userId,
      required this.name,
      this.breed,
      this.age,
      this.photoUrl,
      this.bio,
      required this.createdAt});

  factory _$DogImpl.fromJson(Map<String, dynamic> json) =>
      _$$DogImplFromJson(json);

  @override
  final String id;
  @override
  final String userId;
  @override
  final String name;
  @override
  final String? breed;
  @override
  final int? age;
  @override
  final String? photoUrl;
  @override
  final String? bio;
  @override
  final DateTime createdAt;

  @override
  String toString() {
    return 'Dog(id: $id, userId: $userId, name: $name, breed: $breed, age: $age, photoUrl: $photoUrl, bio: $bio, createdAt: $createdAt)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$DogImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.userId, userId) || other.userId == userId) &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.breed, breed) || other.breed == breed) &&
            (identical(other.age, age) || other.age == age) &&
            (identical(other.photoUrl, photoUrl) ||
                other.photoUrl == photoUrl) &&
            (identical(other.bio, bio) || other.bio == bio) &&
            (identical(other.createdAt, createdAt) ||
                other.createdAt == createdAt));
  }

  @JsonKey(ignore: true)
  @override
  int get hashCode => Object.hash(
      runtimeType, id, userId, name, breed, age, photoUrl, bio, createdAt);

  @JsonKey(ignore: true)
  @override
  @pragma('vm:prefer-inline')
  _$$DogImplCopyWith<_$DogImpl> get copyWith =>
      __$$DogImplCopyWithImpl<_$DogImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$DogImplToJson(
      this,
    );
  }
}

abstract class _Dog implements Dog {
  const factory _Dog(
      {required final String id,
      required final String userId,
      required final String name,
      final String? breed,
      final int? age,
      final String? photoUrl,
      final String? bio,
      required final DateTime createdAt}) = _$DogImpl;

  factory _Dog.fromJson(Map<String, dynamic> json) = _$DogImpl.fromJson;

  @override
  String get id;
  @override
  String get userId;
  @override
  String get name;
  @override
  String? get breed;
  @override
  int? get age;
  @override
  String? get photoUrl;
  @override
  String? get bio;
  @override
  DateTime get createdAt;
  @override
  @JsonKey(ignore: true)
  _$$DogImplCopyWith<_$DogImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

CreateDogRequest _$CreateDogRequestFromJson(Map<String, dynamic> json) {
  return _CreateDogRequest.fromJson(json);
}

/// @nodoc
mixin _$CreateDogRequest {
  String get name => throw _privateConstructorUsedError;
  String? get breed => throw _privateConstructorUsedError;
  int? get age => throw _privateConstructorUsedError;
  String? get photoUrl => throw _privateConstructorUsedError;
  String? get bio => throw _privateConstructorUsedError;

  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  $CreateDogRequestCopyWith<CreateDogRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $CreateDogRequestCopyWith<$Res> {
  factory $CreateDogRequestCopyWith(
          CreateDogRequest value, $Res Function(CreateDogRequest) then) =
      _$CreateDogRequestCopyWithImpl<$Res, CreateDogRequest>;
  @useResult
  $Res call(
      {String name, String? breed, int? age, String? photoUrl, String? bio});
}

/// @nodoc
class _$CreateDogRequestCopyWithImpl<$Res, $Val extends CreateDogRequest>
    implements $CreateDogRequestCopyWith<$Res> {
  _$CreateDogRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = null,
    Object? breed = freezed,
    Object? age = freezed,
    Object? photoUrl = freezed,
    Object? bio = freezed,
  }) {
    return _then(_value.copyWith(
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      breed: freezed == breed
          ? _value.breed
          : breed // ignore: cast_nullable_to_non_nullable
              as String?,
      age: freezed == age
          ? _value.age
          : age // ignore: cast_nullable_to_non_nullable
              as int?,
      photoUrl: freezed == photoUrl
          ? _value.photoUrl
          : photoUrl // ignore: cast_nullable_to_non_nullable
              as String?,
      bio: freezed == bio
          ? _value.bio
          : bio // ignore: cast_nullable_to_non_nullable
              as String?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$CreateDogRequestImplCopyWith<$Res>
    implements $CreateDogRequestCopyWith<$Res> {
  factory _$$CreateDogRequestImplCopyWith(_$CreateDogRequestImpl value,
          $Res Function(_$CreateDogRequestImpl) then) =
      __$$CreateDogRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String name, String? breed, int? age, String? photoUrl, String? bio});
}

/// @nodoc
class __$$CreateDogRequestImplCopyWithImpl<$Res>
    extends _$CreateDogRequestCopyWithImpl<$Res, _$CreateDogRequestImpl>
    implements _$$CreateDogRequestImplCopyWith<$Res> {
  __$$CreateDogRequestImplCopyWithImpl(_$CreateDogRequestImpl _value,
      $Res Function(_$CreateDogRequestImpl) _then)
      : super(_value, _then);

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = null,
    Object? breed = freezed,
    Object? age = freezed,
    Object? photoUrl = freezed,
    Object? bio = freezed,
  }) {
    return _then(_$CreateDogRequestImpl(
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      breed: freezed == breed
          ? _value.breed
          : breed // ignore: cast_nullable_to_non_nullable
              as String?,
      age: freezed == age
          ? _value.age
          : age // ignore: cast_nullable_to_non_nullable
              as int?,
      photoUrl: freezed == photoUrl
          ? _value.photoUrl
          : photoUrl // ignore: cast_nullable_to_non_nullable
              as String?,
      bio: freezed == bio
          ? _value.bio
          : bio // ignore: cast_nullable_to_non_nullable
              as String?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$CreateDogRequestImpl implements _CreateDogRequest {
  const _$CreateDogRequestImpl(
      {required this.name, this.breed, this.age, this.photoUrl, this.bio});

  factory _$CreateDogRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$CreateDogRequestImplFromJson(json);

  @override
  final String name;
  @override
  final String? breed;
  @override
  final int? age;
  @override
  final String? photoUrl;
  @override
  final String? bio;

  @override
  String toString() {
    return 'CreateDogRequest(name: $name, breed: $breed, age: $age, photoUrl: $photoUrl, bio: $bio)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$CreateDogRequestImpl &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.breed, breed) || other.breed == breed) &&
            (identical(other.age, age) || other.age == age) &&
            (identical(other.photoUrl, photoUrl) ||
                other.photoUrl == photoUrl) &&
            (identical(other.bio, bio) || other.bio == bio));
  }

  @JsonKey(ignore: true)
  @override
  int get hashCode => Object.hash(runtimeType, name, breed, age, photoUrl, bio);

  @JsonKey(ignore: true)
  @override
  @pragma('vm:prefer-inline')
  _$$CreateDogRequestImplCopyWith<_$CreateDogRequestImpl> get copyWith =>
      __$$CreateDogRequestImplCopyWithImpl<_$CreateDogRequestImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$CreateDogRequestImplToJson(
      this,
    );
  }
}

abstract class _CreateDogRequest implements CreateDogRequest {
  const factory _CreateDogRequest(
      {required final String name,
      final String? breed,
      final int? age,
      final String? photoUrl,
      final String? bio}) = _$CreateDogRequestImpl;

  factory _CreateDogRequest.fromJson(Map<String, dynamic> json) =
      _$CreateDogRequestImpl.fromJson;

  @override
  String get name;
  @override
  String? get breed;
  @override
  int? get age;
  @override
  String? get photoUrl;
  @override
  String? get bio;
  @override
  @JsonKey(ignore: true)
  _$$CreateDogRequestImplCopyWith<_$CreateDogRequestImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

UpdateDogRequest _$UpdateDogRequestFromJson(Map<String, dynamic> json) {
  return _UpdateDogRequest.fromJson(json);
}

/// @nodoc
mixin _$UpdateDogRequest {
  String? get name => throw _privateConstructorUsedError;
  String? get breed => throw _privateConstructorUsedError;
  int? get age => throw _privateConstructorUsedError;
  String? get photoUrl => throw _privateConstructorUsedError;
  String? get bio => throw _privateConstructorUsedError;

  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  $UpdateDogRequestCopyWith<UpdateDogRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $UpdateDogRequestCopyWith<$Res> {
  factory $UpdateDogRequestCopyWith(
          UpdateDogRequest value, $Res Function(UpdateDogRequest) then) =
      _$UpdateDogRequestCopyWithImpl<$Res, UpdateDogRequest>;
  @useResult
  $Res call(
      {String? name, String? breed, int? age, String? photoUrl, String? bio});
}

/// @nodoc
class _$UpdateDogRequestCopyWithImpl<$Res, $Val extends UpdateDogRequest>
    implements $UpdateDogRequestCopyWith<$Res> {
  _$UpdateDogRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = freezed,
    Object? breed = freezed,
    Object? age = freezed,
    Object? photoUrl = freezed,
    Object? bio = freezed,
  }) {
    return _then(_value.copyWith(
      name: freezed == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String?,
      breed: freezed == breed
          ? _value.breed
          : breed // ignore: cast_nullable_to_non_nullable
              as String?,
      age: freezed == age
          ? _value.age
          : age // ignore: cast_nullable_to_non_nullable
              as int?,
      photoUrl: freezed == photoUrl
          ? _value.photoUrl
          : photoUrl // ignore: cast_nullable_to_non_nullable
              as String?,
      bio: freezed == bio
          ? _value.bio
          : bio // ignore: cast_nullable_to_non_nullable
              as String?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$UpdateDogRequestImplCopyWith<$Res>
    implements $UpdateDogRequestCopyWith<$Res> {
  factory _$$UpdateDogRequestImplCopyWith(_$UpdateDogRequestImpl value,
          $Res Function(_$UpdateDogRequestImpl) then) =
      __$$UpdateDogRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String? name, String? breed, int? age, String? photoUrl, String? bio});
}

/// @nodoc
class __$$UpdateDogRequestImplCopyWithImpl<$Res>
    extends _$UpdateDogRequestCopyWithImpl<$Res, _$UpdateDogRequestImpl>
    implements _$$UpdateDogRequestImplCopyWith<$Res> {
  __$$UpdateDogRequestImplCopyWithImpl(_$UpdateDogRequestImpl _value,
      $Res Function(_$UpdateDogRequestImpl) _then)
      : super(_value, _then);

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = freezed,
    Object? breed = freezed,
    Object? age = freezed,
    Object? photoUrl = freezed,
    Object? bio = freezed,
  }) {
    return _then(_$UpdateDogRequestImpl(
      name: freezed == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String?,
      breed: freezed == breed
          ? _value.breed
          : breed // ignore: cast_nullable_to_non_nullable
              as String?,
      age: freezed == age
          ? _value.age
          : age // ignore: cast_nullable_to_non_nullable
              as int?,
      photoUrl: freezed == photoUrl
          ? _value.photoUrl
          : photoUrl // ignore: cast_nullable_to_non_nullable
              as String?,
      bio: freezed == bio
          ? _value.bio
          : bio // ignore: cast_nullable_to_non_nullable
              as String?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$UpdateDogRequestImpl implements _UpdateDogRequest {
  const _$UpdateDogRequestImpl(
      {this.name, this.breed, this.age, this.photoUrl, this.bio});

  factory _$UpdateDogRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$UpdateDogRequestImplFromJson(json);

  @override
  final String? name;
  @override
  final String? breed;
  @override
  final int? age;
  @override
  final String? photoUrl;
  @override
  final String? bio;

  @override
  String toString() {
    return 'UpdateDogRequest(name: $name, breed: $breed, age: $age, photoUrl: $photoUrl, bio: $bio)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$UpdateDogRequestImpl &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.breed, breed) || other.breed == breed) &&
            (identical(other.age, age) || other.age == age) &&
            (identical(other.photoUrl, photoUrl) ||
                other.photoUrl == photoUrl) &&
            (identical(other.bio, bio) || other.bio == bio));
  }

  @JsonKey(ignore: true)
  @override
  int get hashCode => Object.hash(runtimeType, name, breed, age, photoUrl, bio);

  @JsonKey(ignore: true)
  @override
  @pragma('vm:prefer-inline')
  _$$UpdateDogRequestImplCopyWith<_$UpdateDogRequestImpl> get copyWith =>
      __$$UpdateDogRequestImplCopyWithImpl<_$UpdateDogRequestImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$UpdateDogRequestImplToJson(
      this,
    );
  }
}

abstract class _UpdateDogRequest implements UpdateDogRequest {
  const factory _UpdateDogRequest(
      {final String? name,
      final String? breed,
      final int? age,
      final String? photoUrl,
      final String? bio}) = _$UpdateDogRequestImpl;

  factory _UpdateDogRequest.fromJson(Map<String, dynamic> json) =
      _$UpdateDogRequestImpl.fromJson;

  @override
  String? get name;
  @override
  String? get breed;
  @override
  int? get age;
  @override
  String? get photoUrl;
  @override
  String? get bio;
  @override
  @JsonKey(ignore: true)
  _$$UpdateDogRequestImplCopyWith<_$UpdateDogRequestImpl> get copyWith =>
      throw _privateConstructorUsedError;
}
