import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:doggyclub/models/dog.dart';
import 'package:doggyclub/models/user.dart';

part 'post.freezed.dart';
part 'post.g.dart';

@freezed
class Post with _$Post {
  const factory Post({
    required String id,
    required String dogId,
    required String content,
    @Default([]) List<String> mediaUrls,
    required String mediaType,
    @Default([]) List<String> hashtags,
    String? location,
    @Default(0) int likesCount,
    @Default(0) int commentsCount,
    @Default(0) int sharesCount,
    @Default(0) int viewsCount,
    @Default(true) bool isPublic,
    required DateTime createdAt,
    required DateTime updatedAt,
    required Dog dog,
  }) = _Post;

  factory Post.fromJson(Map<String, dynamic> json) => _$PostFromJson(json);
}

@freezed
class Comment with _$Comment {
  const factory Comment({
    required String id,
    required String postId,
    required String userId,
    required String content,
    @Default(0) int likesCount,
    String? parentId,
    required DateTime createdAt,
    required DateTime updatedAt,
    required User user,
    @Default([]) List<Comment> replies,
  }) = _Comment;

  factory Comment.fromJson(Map<String, dynamic> json) =>
      _$CommentFromJson(json);
}

@freezed
class Like with _$Like {
  const factory Like({
    required String id,
    required String postId,
    required String userId,
    required DateTime createdAt,
  }) = _Like;

  factory Like.fromJson(Map<String, dynamic> json) => _$LikeFromJson(json);
}

@freezed
class Follow with _$Follow {
  const factory Follow({
    required String id,
    required String followerId,
    required String dogId,
    required DateTime createdAt,
  }) = _Follow;

  factory Follow.fromJson(Map<String, dynamic> json) => _$FollowFromJson(json);
}

// Media types
class MediaType {
  static const String photo = 'photo';
  static const String video = 'video';
  static const String mixed = 'mixed';
}