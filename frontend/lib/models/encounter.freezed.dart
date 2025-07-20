// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'encounter.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models');

Encounter _$EncounterFromJson(Map<String, dynamic> json) {
  return _Encounter.fromJson(json);
}

/// @nodoc
mixin _$Encounter {
  String get id => throw _privateConstructorUsedError;
  String get dog1Id => throw _privateConstructorUsedError;
  String get dog2Id => throw _privateConstructorUsedError;
  List<double> get location =>
      throw _privateConstructorUsedError; // [longitude, latitude]
  DateTime get timestamp => throw _privateConstructorUsedError;

  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  $EncounterCopyWith<Encounter> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $EncounterCopyWith<$Res> {
  factory $EncounterCopyWith(Encounter value, $Res Function(Encounter) then) =
      _$EncounterCopyWithImpl<$Res, Encounter>;
  @useResult
  $Res call(
      {String id,
      String dog1Id,
      String dog2Id,
      List<double> location,
      DateTime timestamp});
}

/// @nodoc
class _$EncounterCopyWithImpl<$Res, $Val extends Encounter>
    implements $EncounterCopyWith<$Res> {
  _$EncounterCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? dog1Id = null,
    Object? dog2Id = null,
    Object? location = null,
    Object? timestamp = null,
  }) {
    return _then(_value.copyWith(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as String,
      dog1Id: null == dog1Id
          ? _value.dog1Id
          : dog1Id // ignore: cast_nullable_to_non_nullable
              as String,
      dog2Id: null == dog2Id
          ? _value.dog2Id
          : dog2Id // ignore: cast_nullable_to_non_nullable
              as String,
      location: null == location
          ? _value.location
          : location // ignore: cast_nullable_to_non_nullable
              as List<double>,
      timestamp: null == timestamp
          ? _value.timestamp
          : timestamp // ignore: cast_nullable_to_non_nullable
              as DateTime,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$EncounterImplCopyWith<$Res>
    implements $EncounterCopyWith<$Res> {
  factory _$$EncounterImplCopyWith(
          _$EncounterImpl value, $Res Function(_$EncounterImpl) then) =
      __$$EncounterImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String id,
      String dog1Id,
      String dog2Id,
      List<double> location,
      DateTime timestamp});
}

/// @nodoc
class __$$EncounterImplCopyWithImpl<$Res>
    extends _$EncounterCopyWithImpl<$Res, _$EncounterImpl>
    implements _$$EncounterImplCopyWith<$Res> {
  __$$EncounterImplCopyWithImpl(
      _$EncounterImpl _value, $Res Function(_$EncounterImpl) _then)
      : super(_value, _then);

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? dog1Id = null,
    Object? dog2Id = null,
    Object? location = null,
    Object? timestamp = null,
  }) {
    return _then(_$EncounterImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as String,
      dog1Id: null == dog1Id
          ? _value.dog1Id
          : dog1Id // ignore: cast_nullable_to_non_nullable
              as String,
      dog2Id: null == dog2Id
          ? _value.dog2Id
          : dog2Id // ignore: cast_nullable_to_non_nullable
              as String,
      location: null == location
          ? _value._location
          : location // ignore: cast_nullable_to_non_nullable
              as List<double>,
      timestamp: null == timestamp
          ? _value.timestamp
          : timestamp // ignore: cast_nullable_to_non_nullable
              as DateTime,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$EncounterImpl implements _Encounter {
  const _$EncounterImpl(
      {required this.id,
      required this.dog1Id,
      required this.dog2Id,
      required final List<double> location,
      required this.timestamp})
      : _location = location;

  factory _$EncounterImpl.fromJson(Map<String, dynamic> json) =>
      _$$EncounterImplFromJson(json);

  @override
  final String id;
  @override
  final String dog1Id;
  @override
  final String dog2Id;
  final List<double> _location;
  @override
  List<double> get location {
    if (_location is EqualUnmodifiableListView) return _location;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(_location);
  }

// [longitude, latitude]
  @override
  final DateTime timestamp;

  @override
  String toString() {
    return 'Encounter(id: $id, dog1Id: $dog1Id, dog2Id: $dog2Id, location: $location, timestamp: $timestamp)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$EncounterImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.dog1Id, dog1Id) || other.dog1Id == dog1Id) &&
            (identical(other.dog2Id, dog2Id) || other.dog2Id == dog2Id) &&
            const DeepCollectionEquality().equals(other._location, _location) &&
            (identical(other.timestamp, timestamp) ||
                other.timestamp == timestamp));
  }

  @JsonKey(ignore: true)
  @override
  int get hashCode => Object.hash(runtimeType, id, dog1Id, dog2Id,
      const DeepCollectionEquality().hash(_location), timestamp);

  @JsonKey(ignore: true)
  @override
  @pragma('vm:prefer-inline')
  _$$EncounterImplCopyWith<_$EncounterImpl> get copyWith =>
      __$$EncounterImplCopyWithImpl<_$EncounterImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$EncounterImplToJson(
      this,
    );
  }
}

abstract class _Encounter implements Encounter {
  const factory _Encounter(
      {required final String id,
      required final String dog1Id,
      required final String dog2Id,
      required final List<double> location,
      required final DateTime timestamp}) = _$EncounterImpl;

  factory _Encounter.fromJson(Map<String, dynamic> json) =
      _$EncounterImpl.fromJson;

  @override
  String get id;
  @override
  String get dog1Id;
  @override
  String get dog2Id;
  @override
  List<double> get location;
  @override // [longitude, latitude]
  DateTime get timestamp;
  @override
  @JsonKey(ignore: true)
  _$$EncounterImplCopyWith<_$EncounterImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

DeviceLocation _$DeviceLocationFromJson(Map<String, dynamic> json) {
  return _DeviceLocation.fromJson(json);
}

/// @nodoc
mixin _$DeviceLocation {
  String get id => throw _privateConstructorUsedError;
  String get dogId => throw _privateConstructorUsedError;
  List<double> get location =>
      throw _privateConstructorUsedError; // [longitude, latitude]
  DateTime get timestamp => throw _privateConstructorUsedError;
  DateTime get updatedAt => throw _privateConstructorUsedError;

  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  $DeviceLocationCopyWith<DeviceLocation> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $DeviceLocationCopyWith<$Res> {
  factory $DeviceLocationCopyWith(
          DeviceLocation value, $Res Function(DeviceLocation) then) =
      _$DeviceLocationCopyWithImpl<$Res, DeviceLocation>;
  @useResult
  $Res call(
      {String id,
      String dogId,
      List<double> location,
      DateTime timestamp,
      DateTime updatedAt});
}

/// @nodoc
class _$DeviceLocationCopyWithImpl<$Res, $Val extends DeviceLocation>
    implements $DeviceLocationCopyWith<$Res> {
  _$DeviceLocationCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? dogId = null,
    Object? location = null,
    Object? timestamp = null,
    Object? updatedAt = null,
  }) {
    return _then(_value.copyWith(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as String,
      dogId: null == dogId
          ? _value.dogId
          : dogId // ignore: cast_nullable_to_non_nullable
              as String,
      location: null == location
          ? _value.location
          : location // ignore: cast_nullable_to_non_nullable
              as List<double>,
      timestamp: null == timestamp
          ? _value.timestamp
          : timestamp // ignore: cast_nullable_to_non_nullable
              as DateTime,
      updatedAt: null == updatedAt
          ? _value.updatedAt
          : updatedAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$DeviceLocationImplCopyWith<$Res>
    implements $DeviceLocationCopyWith<$Res> {
  factory _$$DeviceLocationImplCopyWith(_$DeviceLocationImpl value,
          $Res Function(_$DeviceLocationImpl) then) =
      __$$DeviceLocationImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String id,
      String dogId,
      List<double> location,
      DateTime timestamp,
      DateTime updatedAt});
}

/// @nodoc
class __$$DeviceLocationImplCopyWithImpl<$Res>
    extends _$DeviceLocationCopyWithImpl<$Res, _$DeviceLocationImpl>
    implements _$$DeviceLocationImplCopyWith<$Res> {
  __$$DeviceLocationImplCopyWithImpl(
      _$DeviceLocationImpl _value, $Res Function(_$DeviceLocationImpl) _then)
      : super(_value, _then);

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? dogId = null,
    Object? location = null,
    Object? timestamp = null,
    Object? updatedAt = null,
  }) {
    return _then(_$DeviceLocationImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as String,
      dogId: null == dogId
          ? _value.dogId
          : dogId // ignore: cast_nullable_to_non_nullable
              as String,
      location: null == location
          ? _value._location
          : location // ignore: cast_nullable_to_non_nullable
              as List<double>,
      timestamp: null == timestamp
          ? _value.timestamp
          : timestamp // ignore: cast_nullable_to_non_nullable
              as DateTime,
      updatedAt: null == updatedAt
          ? _value.updatedAt
          : updatedAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$DeviceLocationImpl implements _DeviceLocation {
  const _$DeviceLocationImpl(
      {required this.id,
      required this.dogId,
      required final List<double> location,
      required this.timestamp,
      required this.updatedAt})
      : _location = location;

  factory _$DeviceLocationImpl.fromJson(Map<String, dynamic> json) =>
      _$$DeviceLocationImplFromJson(json);

  @override
  final String id;
  @override
  final String dogId;
  final List<double> _location;
  @override
  List<double> get location {
    if (_location is EqualUnmodifiableListView) return _location;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(_location);
  }

// [longitude, latitude]
  @override
  final DateTime timestamp;
  @override
  final DateTime updatedAt;

  @override
  String toString() {
    return 'DeviceLocation(id: $id, dogId: $dogId, location: $location, timestamp: $timestamp, updatedAt: $updatedAt)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$DeviceLocationImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.dogId, dogId) || other.dogId == dogId) &&
            const DeepCollectionEquality().equals(other._location, _location) &&
            (identical(other.timestamp, timestamp) ||
                other.timestamp == timestamp) &&
            (identical(other.updatedAt, updatedAt) ||
                other.updatedAt == updatedAt));
  }

  @JsonKey(ignore: true)
  @override
  int get hashCode => Object.hash(runtimeType, id, dogId,
      const DeepCollectionEquality().hash(_location), timestamp, updatedAt);

  @JsonKey(ignore: true)
  @override
  @pragma('vm:prefer-inline')
  _$$DeviceLocationImplCopyWith<_$DeviceLocationImpl> get copyWith =>
      __$$DeviceLocationImplCopyWithImpl<_$DeviceLocationImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$DeviceLocationImplToJson(
      this,
    );
  }
}

abstract class _DeviceLocation implements DeviceLocation {
  const factory _DeviceLocation(
      {required final String id,
      required final String dogId,
      required final List<double> location,
      required final DateTime timestamp,
      required final DateTime updatedAt}) = _$DeviceLocationImpl;

  factory _DeviceLocation.fromJson(Map<String, dynamic> json) =
      _$DeviceLocationImpl.fromJson;

  @override
  String get id;
  @override
  String get dogId;
  @override
  List<double> get location;
  @override // [longitude, latitude]
  DateTime get timestamp;
  @override
  DateTime get updatedAt;
  @override
  @JsonKey(ignore: true)
  _$$DeviceLocationImplCopyWith<_$DeviceLocationImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

DetectEncountersRequest _$DetectEncountersRequestFromJson(
    Map<String, dynamic> json) {
  return _DetectEncountersRequest.fromJson(json);
}

/// @nodoc
mixin _$DetectEncountersRequest {
  String get dogId => throw _privateConstructorUsedError;
  double get radiusMeters => throw _privateConstructorUsedError;

  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  $DetectEncountersRequestCopyWith<DetectEncountersRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $DetectEncountersRequestCopyWith<$Res> {
  factory $DetectEncountersRequestCopyWith(DetectEncountersRequest value,
          $Res Function(DetectEncountersRequest) then) =
      _$DetectEncountersRequestCopyWithImpl<$Res, DetectEncountersRequest>;
  @useResult
  $Res call({String dogId, double radiusMeters});
}

/// @nodoc
class _$DetectEncountersRequestCopyWithImpl<$Res,
        $Val extends DetectEncountersRequest>
    implements $DetectEncountersRequestCopyWith<$Res> {
  _$DetectEncountersRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? dogId = null,
    Object? radiusMeters = null,
  }) {
    return _then(_value.copyWith(
      dogId: null == dogId
          ? _value.dogId
          : dogId // ignore: cast_nullable_to_non_nullable
              as String,
      radiusMeters: null == radiusMeters
          ? _value.radiusMeters
          : radiusMeters // ignore: cast_nullable_to_non_nullable
              as double,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$DetectEncountersRequestImplCopyWith<$Res>
    implements $DetectEncountersRequestCopyWith<$Res> {
  factory _$$DetectEncountersRequestImplCopyWith(
          _$DetectEncountersRequestImpl value,
          $Res Function(_$DetectEncountersRequestImpl) then) =
      __$$DetectEncountersRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({String dogId, double radiusMeters});
}

/// @nodoc
class __$$DetectEncountersRequestImplCopyWithImpl<$Res>
    extends _$DetectEncountersRequestCopyWithImpl<$Res,
        _$DetectEncountersRequestImpl>
    implements _$$DetectEncountersRequestImplCopyWith<$Res> {
  __$$DetectEncountersRequestImplCopyWithImpl(
      _$DetectEncountersRequestImpl _value,
      $Res Function(_$DetectEncountersRequestImpl) _then)
      : super(_value, _then);

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? dogId = null,
    Object? radiusMeters = null,
  }) {
    return _then(_$DetectEncountersRequestImpl(
      dogId: null == dogId
          ? _value.dogId
          : dogId // ignore: cast_nullable_to_non_nullable
              as String,
      radiusMeters: null == radiusMeters
          ? _value.radiusMeters
          : radiusMeters // ignore: cast_nullable_to_non_nullable
              as double,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$DetectEncountersRequestImpl implements _DetectEncountersRequest {
  const _$DetectEncountersRequestImpl(
      {required this.dogId, required this.radiusMeters});

  factory _$DetectEncountersRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$DetectEncountersRequestImplFromJson(json);

  @override
  final String dogId;
  @override
  final double radiusMeters;

  @override
  String toString() {
    return 'DetectEncountersRequest(dogId: $dogId, radiusMeters: $radiusMeters)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$DetectEncountersRequestImpl &&
            (identical(other.dogId, dogId) || other.dogId == dogId) &&
            (identical(other.radiusMeters, radiusMeters) ||
                other.radiusMeters == radiusMeters));
  }

  @JsonKey(ignore: true)
  @override
  int get hashCode => Object.hash(runtimeType, dogId, radiusMeters);

  @JsonKey(ignore: true)
  @override
  @pragma('vm:prefer-inline')
  _$$DetectEncountersRequestImplCopyWith<_$DetectEncountersRequestImpl>
      get copyWith => __$$DetectEncountersRequestImplCopyWithImpl<
          _$DetectEncountersRequestImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$DetectEncountersRequestImplToJson(
      this,
    );
  }
}

abstract class _DetectEncountersRequest implements DetectEncountersRequest {
  const factory _DetectEncountersRequest(
      {required final String dogId,
      required final double radiusMeters}) = _$DetectEncountersRequestImpl;

  factory _DetectEncountersRequest.fromJson(Map<String, dynamic> json) =
      _$DetectEncountersRequestImpl.fromJson;

  @override
  String get dogId;
  @override
  double get radiusMeters;
  @override
  @JsonKey(ignore: true)
  _$$DetectEncountersRequestImplCopyWith<_$DetectEncountersRequestImpl>
      get copyWith => throw _privateConstructorUsedError;
}

EncounterResponse _$EncounterResponseFromJson(Map<String, dynamic> json) {
  return _EncounterResponse.fromJson(json);
}

/// @nodoc
mixin _$EncounterResponse {
  List<Encounter> get encounters => throw _privateConstructorUsedError;
  int get total => throw _privateConstructorUsedError;
  int get limit => throw _privateConstructorUsedError;
  int get offset => throw _privateConstructorUsedError;

  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  $EncounterResponseCopyWith<EncounterResponse> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $EncounterResponseCopyWith<$Res> {
  factory $EncounterResponseCopyWith(
          EncounterResponse value, $Res Function(EncounterResponse) then) =
      _$EncounterResponseCopyWithImpl<$Res, EncounterResponse>;
  @useResult
  $Res call({List<Encounter> encounters, int total, int limit, int offset});
}

/// @nodoc
class _$EncounterResponseCopyWithImpl<$Res, $Val extends EncounterResponse>
    implements $EncounterResponseCopyWith<$Res> {
  _$EncounterResponseCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? encounters = null,
    Object? total = null,
    Object? limit = null,
    Object? offset = null,
  }) {
    return _then(_value.copyWith(
      encounters: null == encounters
          ? _value.encounters
          : encounters // ignore: cast_nullable_to_non_nullable
              as List<Encounter>,
      total: null == total
          ? _value.total
          : total // ignore: cast_nullable_to_non_nullable
              as int,
      limit: null == limit
          ? _value.limit
          : limit // ignore: cast_nullable_to_non_nullable
              as int,
      offset: null == offset
          ? _value.offset
          : offset // ignore: cast_nullable_to_non_nullable
              as int,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$EncounterResponseImplCopyWith<$Res>
    implements $EncounterResponseCopyWith<$Res> {
  factory _$$EncounterResponseImplCopyWith(_$EncounterResponseImpl value,
          $Res Function(_$EncounterResponseImpl) then) =
      __$$EncounterResponseImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({List<Encounter> encounters, int total, int limit, int offset});
}

/// @nodoc
class __$$EncounterResponseImplCopyWithImpl<$Res>
    extends _$EncounterResponseCopyWithImpl<$Res, _$EncounterResponseImpl>
    implements _$$EncounterResponseImplCopyWith<$Res> {
  __$$EncounterResponseImplCopyWithImpl(_$EncounterResponseImpl _value,
      $Res Function(_$EncounterResponseImpl) _then)
      : super(_value, _then);

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? encounters = null,
    Object? total = null,
    Object? limit = null,
    Object? offset = null,
  }) {
    return _then(_$EncounterResponseImpl(
      encounters: null == encounters
          ? _value._encounters
          : encounters // ignore: cast_nullable_to_non_nullable
              as List<Encounter>,
      total: null == total
          ? _value.total
          : total // ignore: cast_nullable_to_non_nullable
              as int,
      limit: null == limit
          ? _value.limit
          : limit // ignore: cast_nullable_to_non_nullable
              as int,
      offset: null == offset
          ? _value.offset
          : offset // ignore: cast_nullable_to_non_nullable
              as int,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$EncounterResponseImpl implements _EncounterResponse {
  const _$EncounterResponseImpl(
      {required final List<Encounter> encounters,
      required this.total,
      required this.limit,
      required this.offset})
      : _encounters = encounters;

  factory _$EncounterResponseImpl.fromJson(Map<String, dynamic> json) =>
      _$$EncounterResponseImplFromJson(json);

  final List<Encounter> _encounters;
  @override
  List<Encounter> get encounters {
    if (_encounters is EqualUnmodifiableListView) return _encounters;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(_encounters);
  }

  @override
  final int total;
  @override
  final int limit;
  @override
  final int offset;

  @override
  String toString() {
    return 'EncounterResponse(encounters: $encounters, total: $total, limit: $limit, offset: $offset)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$EncounterResponseImpl &&
            const DeepCollectionEquality()
                .equals(other._encounters, _encounters) &&
            (identical(other.total, total) || other.total == total) &&
            (identical(other.limit, limit) || other.limit == limit) &&
            (identical(other.offset, offset) || other.offset == offset));
  }

  @JsonKey(ignore: true)
  @override
  int get hashCode => Object.hash(runtimeType,
      const DeepCollectionEquality().hash(_encounters), total, limit, offset);

  @JsonKey(ignore: true)
  @override
  @pragma('vm:prefer-inline')
  _$$EncounterResponseImplCopyWith<_$EncounterResponseImpl> get copyWith =>
      __$$EncounterResponseImplCopyWithImpl<_$EncounterResponseImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$EncounterResponseImplToJson(
      this,
    );
  }
}

abstract class _EncounterResponse implements EncounterResponse {
  const factory _EncounterResponse(
      {required final List<Encounter> encounters,
      required final int total,
      required final int limit,
      required final int offset}) = _$EncounterResponseImpl;

  factory _EncounterResponse.fromJson(Map<String, dynamic> json) =
      _$EncounterResponseImpl.fromJson;

  @override
  List<Encounter> get encounters;
  @override
  int get total;
  @override
  int get limit;
  @override
  int get offset;
  @override
  @JsonKey(ignore: true)
  _$$EncounterResponseImplCopyWith<_$EncounterResponseImpl> get copyWith =>
      throw _privateConstructorUsedError;
}
