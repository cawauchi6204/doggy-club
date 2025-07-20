import 'dart:async';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:doggyclub/models/notification.dart';
import 'package:doggyclub/services/notification_service.dart';
import 'package:doggyclub/providers/auth_provider.dart';

// Notification state
class NotificationState {
  final List<AppNotification> notifications;
  final int unreadCount;
  final bool isLoading;
  final String? error;
  final bool hasPermission;
  final String? fcmToken;

  const NotificationState({
    this.notifications = const [],
    this.unreadCount = 0,
    this.isLoading = false,
    this.error,
    this.hasPermission = false,
    this.fcmToken,
  });

  NotificationState copyWith({
    List<AppNotification>? notifications,
    int? unreadCount,
    bool? isLoading,
    String? error,
    bool? hasPermission,
    String? fcmToken,
  }) {
    return NotificationState(
      notifications: notifications ?? this.notifications,
      unreadCount: unreadCount ?? this.unreadCount,
      isLoading: isLoading ?? this.isLoading,
      error: error,
      hasPermission: hasPermission ?? this.hasPermission,
      fcmToken: fcmToken ?? this.fcmToken,
    );
  }
}

// Notification notifier
class NotificationNotifier extends StateNotifier<NotificationState> {
  final NotificationService _notificationService;
  final Ref _ref;

  NotificationNotifier(
    this._notificationService,
    this._ref,
  ) : super(const NotificationState());

  Future<void> loadNotifications({
    int limit = 20,
    int offset = 0,
  }) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final authState = _ref.read(authProvider);
      if (authState.token == null || authState.user == null) {
        state = state.copyWith(
          error: 'User not authenticated',
          isLoading: false,
        );
        return;
      }

      final response = await _notificationService.getNotifications(
        limit: limit,
        offset: offset,
      );

      List<AppNotification> allNotifications;
      if (offset == 0) {
        allNotifications = response.notifications;
      } else {
        allNotifications = [...state.notifications, ...response.notifications];
      }

      state = state.copyWith(
        notifications: allNotifications,
        isLoading: false,
      );
    } catch (e) {
      print('Error loading notifications: $e');
      state = state.copyWith(
        error: e.toString(),
        isLoading: false,
      );
    }
  }

  Future<void> registerDevice(String deviceToken) async {
    try {
      final request = RegisterDeviceRequest(
        deviceToken: deviceToken,
        platform: 'android', // Simplified for now
      );

      await _notificationService.registerDevice(request);
      state = state.copyWith(fcmToken: deviceToken);
    } catch (e) {
      print('Error registering device: $e');
      state = state.copyWith(error: e.toString());
    }
  }

  void clearError() {
    state = state.copyWith(error: null);
  }

  @override
  void dispose() {
    super.dispose();
  }
}

// Providers
final notificationProvider = StateNotifierProvider<NotificationNotifier, NotificationState>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return NotificationNotifier(
    apiClient.notificationService,
    ref,
  );
});

// Convenience providers
final notificationsProvider = Provider<List<AppNotification>>((ref) {
  return ref.watch(notificationProvider).notifications;
});

final unreadCountProvider = Provider<int>((ref) {
  return ref.watch(notificationProvider).unreadCount;
});

final hasNotificationPermissionProvider = Provider<bool>((ref) {
  return ref.watch(notificationProvider).hasPermission;
});