// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'gift.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models');

Gift _$GiftFromJson(Map<String, dynamic> json) {
  return _Gift.fromJson(json);
}

/// @nodoc
mixin _$Gift {
  String get id => throw _privateConstructorUsedError;
  String get senderDogId => throw _privateConstructorUsedError;
  String get receiverDogId => throw _privateConstructorUsedError;
  String get type => throw _privateConstructorUsedError;
  DateTime get createdAt => throw _privateConstructorUsedError;

  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  $GiftCopyWith<Gift> get copyWith => throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $GiftCopyWith<$Res> {
  factory $GiftCopyWith(Gift value, $Res Function(Gift) then) =
      _$GiftCopyWithImpl<$Res, Gift>;
  @useResult
  $Res call(
      {String id,
      String senderDogId,
      String receiverDogId,
      String type,
      DateTime createdAt});
}

/// @nodoc
class _$GiftCopyWithImpl<$Res, $Val extends Gift>
    implements $GiftCopyWith<$Res> {
  _$GiftCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? senderDogId = null,
    Object? receiverDogId = null,
    Object? type = null,
    Object? createdAt = null,
  }) {
    return _then(_value.copyWith(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as String,
      senderDogId: null == senderDogId
          ? _value.senderDogId
          : senderDogId // ignore: cast_nullable_to_non_nullable
              as String,
      receiverDogId: null == receiverDogId
          ? _value.receiverDogId
          : receiverDogId // ignore: cast_nullable_to_non_nullable
              as String,
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as String,
      createdAt: null == createdAt
          ? _value.createdAt
          : createdAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$GiftImplCopyWith<$Res> implements $GiftCopyWith<$Res> {
  factory _$$GiftImplCopyWith(
          _$GiftImpl value, $Res Function(_$GiftImpl) then) =
      __$$GiftImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String id,
      String senderDogId,
      String receiverDogId,
      String type,
      DateTime createdAt});
}

/// @nodoc
class __$$GiftImplCopyWithImpl<$Res>
    extends _$GiftCopyWithImpl<$Res, _$GiftImpl>
    implements _$$GiftImplCopyWith<$Res> {
  __$$GiftImplCopyWithImpl(_$GiftImpl _value, $Res Function(_$GiftImpl) _then)
      : super(_value, _then);

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? senderDogId = null,
    Object? receiverDogId = null,
    Object? type = null,
    Object? createdAt = null,
  }) {
    return _then(_$GiftImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as String,
      senderDogId: null == senderDogId
          ? _value.senderDogId
          : senderDogId // ignore: cast_nullable_to_non_nullable
              as String,
      receiverDogId: null == receiverDogId
          ? _value.receiverDogId
          : receiverDogId // ignore: cast_nullable_to_non_nullable
              as String,
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as String,
      createdAt: null == createdAt
          ? _value.createdAt
          : createdAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$GiftImpl implements _Gift {
  const _$GiftImpl(
      {required this.id,
      required this.senderDogId,
      required this.receiverDogId,
      required this.type,
      required this.createdAt});

  factory _$GiftImpl.fromJson(Map<String, dynamic> json) =>
      _$$GiftImplFromJson(json);

  @override
  final String id;
  @override
  final String senderDogId;
  @override
  final String receiverDogId;
  @override
  final String type;
  @override
  final DateTime createdAt;

  @override
  String toString() {
    return 'Gift(id: $id, senderDogId: $senderDogId, receiverDogId: $receiverDogId, type: $type, createdAt: $createdAt)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$GiftImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.senderDogId, senderDogId) ||
                other.senderDogId == senderDogId) &&
            (identical(other.receiverDogId, receiverDogId) ||
                other.receiverDogId == receiverDogId) &&
            (identical(other.type, type) || other.type == type) &&
            (identical(other.createdAt, createdAt) ||
                other.createdAt == createdAt));
  }

  @JsonKey(ignore: true)
  @override
  int get hashCode =>
      Object.hash(runtimeType, id, senderDogId, receiverDogId, type, createdAt);

  @JsonKey(ignore: true)
  @override
  @pragma('vm:prefer-inline')
  _$$GiftImplCopyWith<_$GiftImpl> get copyWith =>
      __$$GiftImplCopyWithImpl<_$GiftImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$GiftImplToJson(
      this,
    );
  }
}

abstract class _Gift implements Gift {
  const factory _Gift(
      {required final String id,
      required final String senderDogId,
      required final String receiverDogId,
      required final String type,
      required final DateTime createdAt}) = _$GiftImpl;

  factory _Gift.fromJson(Map<String, dynamic> json) = _$GiftImpl.fromJson;

  @override
  String get id;
  @override
  String get senderDogId;
  @override
  String get receiverDogId;
  @override
  String get type;
  @override
  DateTime get createdAt;
  @override
  @JsonKey(ignore: true)
  _$$GiftImplCopyWith<_$GiftImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

SendGiftRequest _$SendGiftRequestFromJson(Map<String, dynamic> json) {
  return _SendGiftRequest.fromJson(json);
}

/// @nodoc
mixin _$SendGiftRequest {
  String get senderDogId => throw _privateConstructorUsedError;
  String get receiverDogId => throw _privateConstructorUsedError;
  String get type => throw _privateConstructorUsedError;

  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  $SendGiftRequestCopyWith<SendGiftRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $SendGiftRequestCopyWith<$Res> {
  factory $SendGiftRequestCopyWith(
          SendGiftRequest value, $Res Function(SendGiftRequest) then) =
      _$SendGiftRequestCopyWithImpl<$Res, SendGiftRequest>;
  @useResult
  $Res call({String senderDogId, String receiverDogId, String type});
}

/// @nodoc
class _$SendGiftRequestCopyWithImpl<$Res, $Val extends SendGiftRequest>
    implements $SendGiftRequestCopyWith<$Res> {
  _$SendGiftRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? senderDogId = null,
    Object? receiverDogId = null,
    Object? type = null,
  }) {
    return _then(_value.copyWith(
      senderDogId: null == senderDogId
          ? _value.senderDogId
          : senderDogId // ignore: cast_nullable_to_non_nullable
              as String,
      receiverDogId: null == receiverDogId
          ? _value.receiverDogId
          : receiverDogId // ignore: cast_nullable_to_non_nullable
              as String,
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as String,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$SendGiftRequestImplCopyWith<$Res>
    implements $SendGiftRequestCopyWith<$Res> {
  factory _$$SendGiftRequestImplCopyWith(_$SendGiftRequestImpl value,
          $Res Function(_$SendGiftRequestImpl) then) =
      __$$SendGiftRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({String senderDogId, String receiverDogId, String type});
}

/// @nodoc
class __$$SendGiftRequestImplCopyWithImpl<$Res>
    extends _$SendGiftRequestCopyWithImpl<$Res, _$SendGiftRequestImpl>
    implements _$$SendGiftRequestImplCopyWith<$Res> {
  __$$SendGiftRequestImplCopyWithImpl(
      _$SendGiftRequestImpl _value, $Res Function(_$SendGiftRequestImpl) _then)
      : super(_value, _then);

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? senderDogId = null,
    Object? receiverDogId = null,
    Object? type = null,
  }) {
    return _then(_$SendGiftRequestImpl(
      senderDogId: null == senderDogId
          ? _value.senderDogId
          : senderDogId // ignore: cast_nullable_to_non_nullable
              as String,
      receiverDogId: null == receiverDogId
          ? _value.receiverDogId
          : receiverDogId // ignore: cast_nullable_to_non_nullable
              as String,
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as String,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$SendGiftRequestImpl implements _SendGiftRequest {
  const _$SendGiftRequestImpl(
      {required this.senderDogId,
      required this.receiverDogId,
      required this.type});

  factory _$SendGiftRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$SendGiftRequestImplFromJson(json);

  @override
  final String senderDogId;
  @override
  final String receiverDogId;
  @override
  final String type;

  @override
  String toString() {
    return 'SendGiftRequest(senderDogId: $senderDogId, receiverDogId: $receiverDogId, type: $type)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$SendGiftRequestImpl &&
            (identical(other.senderDogId, senderDogId) ||
                other.senderDogId == senderDogId) &&
            (identical(other.receiverDogId, receiverDogId) ||
                other.receiverDogId == receiverDogId) &&
            (identical(other.type, type) || other.type == type));
  }

  @JsonKey(ignore: true)
  @override
  int get hashCode =>
      Object.hash(runtimeType, senderDogId, receiverDogId, type);

  @JsonKey(ignore: true)
  @override
  @pragma('vm:prefer-inline')
  _$$SendGiftRequestImplCopyWith<_$SendGiftRequestImpl> get copyWith =>
      __$$SendGiftRequestImplCopyWithImpl<_$SendGiftRequestImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$SendGiftRequestImplToJson(
      this,
    );
  }
}

abstract class _SendGiftRequest implements SendGiftRequest {
  const factory _SendGiftRequest(
      {required final String senderDogId,
      required final String receiverDogId,
      required final String type}) = _$SendGiftRequestImpl;

  factory _SendGiftRequest.fromJson(Map<String, dynamic> json) =
      _$SendGiftRequestImpl.fromJson;

  @override
  String get senderDogId;
  @override
  String get receiverDogId;
  @override
  String get type;
  @override
  @JsonKey(ignore: true)
  _$$SendGiftRequestImplCopyWith<_$SendGiftRequestImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

GiftResponse _$GiftResponseFromJson(Map<String, dynamic> json) {
  return _GiftResponse.fromJson(json);
}

/// @nodoc
mixin _$GiftResponse {
  List<Gift> get gifts => throw _privateConstructorUsedError;
  int get total => throw _privateConstructorUsedError;
  int get limit => throw _privateConstructorUsedError;
  int get offset => throw _privateConstructorUsedError;

  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  $GiftResponseCopyWith<GiftResponse> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $GiftResponseCopyWith<$Res> {
  factory $GiftResponseCopyWith(
          GiftResponse value, $Res Function(GiftResponse) then) =
      _$GiftResponseCopyWithImpl<$Res, GiftResponse>;
  @useResult
  $Res call({List<Gift> gifts, int total, int limit, int offset});
}

/// @nodoc
class _$GiftResponseCopyWithImpl<$Res, $Val extends GiftResponse>
    implements $GiftResponseCopyWith<$Res> {
  _$GiftResponseCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? gifts = null,
    Object? total = null,
    Object? limit = null,
    Object? offset = null,
  }) {
    return _then(_value.copyWith(
      gifts: null == gifts
          ? _value.gifts
          : gifts // ignore: cast_nullable_to_non_nullable
              as List<Gift>,
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
abstract class _$$GiftResponseImplCopyWith<$Res>
    implements $GiftResponseCopyWith<$Res> {
  factory _$$GiftResponseImplCopyWith(
          _$GiftResponseImpl value, $Res Function(_$GiftResponseImpl) then) =
      __$$GiftResponseImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({List<Gift> gifts, int total, int limit, int offset});
}

/// @nodoc
class __$$GiftResponseImplCopyWithImpl<$Res>
    extends _$GiftResponseCopyWithImpl<$Res, _$GiftResponseImpl>
    implements _$$GiftResponseImplCopyWith<$Res> {
  __$$GiftResponseImplCopyWithImpl(
      _$GiftResponseImpl _value, $Res Function(_$GiftResponseImpl) _then)
      : super(_value, _then);

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? gifts = null,
    Object? total = null,
    Object? limit = null,
    Object? offset = null,
  }) {
    return _then(_$GiftResponseImpl(
      gifts: null == gifts
          ? _value._gifts
          : gifts // ignore: cast_nullable_to_non_nullable
              as List<Gift>,
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
class _$GiftResponseImpl implements _GiftResponse {
  const _$GiftResponseImpl(
      {required final List<Gift> gifts,
      required this.total,
      required this.limit,
      required this.offset})
      : _gifts = gifts;

  factory _$GiftResponseImpl.fromJson(Map<String, dynamic> json) =>
      _$$GiftResponseImplFromJson(json);

  final List<Gift> _gifts;
  @override
  List<Gift> get gifts {
    if (_gifts is EqualUnmodifiableListView) return _gifts;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(_gifts);
  }

  @override
  final int total;
  @override
  final int limit;
  @override
  final int offset;

  @override
  String toString() {
    return 'GiftResponse(gifts: $gifts, total: $total, limit: $limit, offset: $offset)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$GiftResponseImpl &&
            const DeepCollectionEquality().equals(other._gifts, _gifts) &&
            (identical(other.total, total) || other.total == total) &&
            (identical(other.limit, limit) || other.limit == limit) &&
            (identical(other.offset, offset) || other.offset == offset));
  }

  @JsonKey(ignore: true)
  @override
  int get hashCode => Object.hash(runtimeType,
      const DeepCollectionEquality().hash(_gifts), total, limit, offset);

  @JsonKey(ignore: true)
  @override
  @pragma('vm:prefer-inline')
  _$$GiftResponseImplCopyWith<_$GiftResponseImpl> get copyWith =>
      __$$GiftResponseImplCopyWithImpl<_$GiftResponseImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$GiftResponseImplToJson(
      this,
    );
  }
}

abstract class _GiftResponse implements GiftResponse {
  const factory _GiftResponse(
      {required final List<Gift> gifts,
      required final int total,
      required final int limit,
      required final int offset}) = _$GiftResponseImpl;

  factory _GiftResponse.fromJson(Map<String, dynamic> json) =
      _$GiftResponseImpl.fromJson;

  @override
  List<Gift> get gifts;
  @override
  int get total;
  @override
  int get limit;
  @override
  int get offset;
  @override
  @JsonKey(ignore: true)
  _$$GiftResponseImplCopyWith<_$GiftResponseImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

GiftTypeInfo _$GiftTypeInfoFromJson(Map<String, dynamic> json) {
  return _GiftTypeInfo.fromJson(json);
}

/// @nodoc
mixin _$GiftTypeInfo {
  String get type => throw _privateConstructorUsedError;
  String get name => throw _privateConstructorUsedError;

  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  $GiftTypeInfoCopyWith<GiftTypeInfo> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $GiftTypeInfoCopyWith<$Res> {
  factory $GiftTypeInfoCopyWith(
          GiftTypeInfo value, $Res Function(GiftTypeInfo) then) =
      _$GiftTypeInfoCopyWithImpl<$Res, GiftTypeInfo>;
  @useResult
  $Res call({String type, String name});
}

/// @nodoc
class _$GiftTypeInfoCopyWithImpl<$Res, $Val extends GiftTypeInfo>
    implements $GiftTypeInfoCopyWith<$Res> {
  _$GiftTypeInfoCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? type = null,
    Object? name = null,
  }) {
    return _then(_value.copyWith(
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as String,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$GiftTypeInfoImplCopyWith<$Res>
    implements $GiftTypeInfoCopyWith<$Res> {
  factory _$$GiftTypeInfoImplCopyWith(
          _$GiftTypeInfoImpl value, $Res Function(_$GiftTypeInfoImpl) then) =
      __$$GiftTypeInfoImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({String type, String name});
}

/// @nodoc
class __$$GiftTypeInfoImplCopyWithImpl<$Res>
    extends _$GiftTypeInfoCopyWithImpl<$Res, _$GiftTypeInfoImpl>
    implements _$$GiftTypeInfoImplCopyWith<$Res> {
  __$$GiftTypeInfoImplCopyWithImpl(
      _$GiftTypeInfoImpl _value, $Res Function(_$GiftTypeInfoImpl) _then)
      : super(_value, _then);

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? type = null,
    Object? name = null,
  }) {
    return _then(_$GiftTypeInfoImpl(
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as String,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$GiftTypeInfoImpl implements _GiftTypeInfo {
  const _$GiftTypeInfoImpl({required this.type, required this.name});

  factory _$GiftTypeInfoImpl.fromJson(Map<String, dynamic> json) =>
      _$$GiftTypeInfoImplFromJson(json);

  @override
  final String type;
  @override
  final String name;

  @override
  String toString() {
    return 'GiftTypeInfo(type: $type, name: $name)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$GiftTypeInfoImpl &&
            (identical(other.type, type) || other.type == type) &&
            (identical(other.name, name) || other.name == name));
  }

  @JsonKey(ignore: true)
  @override
  int get hashCode => Object.hash(runtimeType, type, name);

  @JsonKey(ignore: true)
  @override
  @pragma('vm:prefer-inline')
  _$$GiftTypeInfoImplCopyWith<_$GiftTypeInfoImpl> get copyWith =>
      __$$GiftTypeInfoImplCopyWithImpl<_$GiftTypeInfoImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$GiftTypeInfoImplToJson(
      this,
    );
  }
}

abstract class _GiftTypeInfo implements GiftTypeInfo {
  const factory _GiftTypeInfo(
      {required final String type,
      required final String name}) = _$GiftTypeInfoImpl;

  factory _GiftTypeInfo.fromJson(Map<String, dynamic> json) =
      _$GiftTypeInfoImpl.fromJson;

  @override
  String get type;
  @override
  String get name;
  @override
  @JsonKey(ignore: true)
  _$$GiftTypeInfoImplCopyWith<_$GiftTypeInfoImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

GiftRanking _$GiftRankingFromJson(Map<String, dynamic> json) {
  return _GiftRanking.fromJson(json);
}

/// @nodoc
mixin _$GiftRanking {
  String get type => throw _privateConstructorUsedError;
  int get count => throw _privateConstructorUsedError;

  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  $GiftRankingCopyWith<GiftRanking> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $GiftRankingCopyWith<$Res> {
  factory $GiftRankingCopyWith(
          GiftRanking value, $Res Function(GiftRanking) then) =
      _$GiftRankingCopyWithImpl<$Res, GiftRanking>;
  @useResult
  $Res call({String type, int count});
}

/// @nodoc
class _$GiftRankingCopyWithImpl<$Res, $Val extends GiftRanking>
    implements $GiftRankingCopyWith<$Res> {
  _$GiftRankingCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? type = null,
    Object? count = null,
  }) {
    return _then(_value.copyWith(
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as String,
      count: null == count
          ? _value.count
          : count // ignore: cast_nullable_to_non_nullable
              as int,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$GiftRankingImplCopyWith<$Res>
    implements $GiftRankingCopyWith<$Res> {
  factory _$$GiftRankingImplCopyWith(
          _$GiftRankingImpl value, $Res Function(_$GiftRankingImpl) then) =
      __$$GiftRankingImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({String type, int count});
}

/// @nodoc
class __$$GiftRankingImplCopyWithImpl<$Res>
    extends _$GiftRankingCopyWithImpl<$Res, _$GiftRankingImpl>
    implements _$$GiftRankingImplCopyWith<$Res> {
  __$$GiftRankingImplCopyWithImpl(
      _$GiftRankingImpl _value, $Res Function(_$GiftRankingImpl) _then)
      : super(_value, _then);

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? type = null,
    Object? count = null,
  }) {
    return _then(_$GiftRankingImpl(
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as String,
      count: null == count
          ? _value.count
          : count // ignore: cast_nullable_to_non_nullable
              as int,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$GiftRankingImpl implements _GiftRanking {
  const _$GiftRankingImpl({required this.type, required this.count});

  factory _$GiftRankingImpl.fromJson(Map<String, dynamic> json) =>
      _$$GiftRankingImplFromJson(json);

  @override
  final String type;
  @override
  final int count;

  @override
  String toString() {
    return 'GiftRanking(type: $type, count: $count)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$GiftRankingImpl &&
            (identical(other.type, type) || other.type == type) &&
            (identical(other.count, count) || other.count == count));
  }

  @JsonKey(ignore: true)
  @override
  int get hashCode => Object.hash(runtimeType, type, count);

  @JsonKey(ignore: true)
  @override
  @pragma('vm:prefer-inline')
  _$$GiftRankingImplCopyWith<_$GiftRankingImpl> get copyWith =>
      __$$GiftRankingImplCopyWithImpl<_$GiftRankingImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$GiftRankingImplToJson(
      this,
    );
  }
}

abstract class _GiftRanking implements GiftRanking {
  const factory _GiftRanking(
      {required final String type,
      required final int count}) = _$GiftRankingImpl;

  factory _GiftRanking.fromJson(Map<String, dynamic> json) =
      _$GiftRankingImpl.fromJson;

  @override
  String get type;
  @override
  int get count;
  @override
  @JsonKey(ignore: true)
  _$$GiftRankingImplCopyWith<_$GiftRankingImpl> get copyWith =>
      throw _privateConstructorUsedError;
}
