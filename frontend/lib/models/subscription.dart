import 'package:freezed_annotation/freezed_annotation.dart';

part 'subscription.freezed.dart';
part 'subscription.g.dart';

@freezed
class Subscription with _$Subscription {
  const factory Subscription({
    required String id,
    required String userId,
    required String type,
    required String status,
    required String stripeCustomerId,
    required String stripeSubscriptionId,
    required DateTime currentPeriodStart,
    required DateTime currentPeriodEnd,
    @Default(false) bool cancelAtPeriodEnd,
    required DateTime createdAt,
    required DateTime updatedAt,
  }) = _Subscription;

  factory Subscription.fromJson(Map<String, dynamic> json) =>
      _$SubscriptionFromJson(json);
}

// Constants
class SubscriptionType {
  static const String premiumMonthly = 'premium_monthly';
  static const String premiumYearly = 'premium_yearly';
}

class SubscriptionStatus {
  static const String active = 'active';
  static const String cancelled = 'cancelled';
  static const String expired = 'expired';
}