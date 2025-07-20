import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:doggyclub/providers/auth_provider.dart';
import 'package:doggyclub/screens/auth/login_screen.dart';
import 'package:doggyclub/screens/auth/register_screen.dart';
import 'package:doggyclub/screens/profile/profile_screen.dart';
import 'package:doggyclub/screens/dogs/add_dog_screen.dart';

void main() {
  runApp(
    const ProviderScope(
      child: MyApp(),
    ),
  );
}

// Alias for testing compatibility
class MyApp extends DoggyClubApp {
  const MyApp({super.key});
}

class DoggyClubApp extends ConsumerWidget {
  const DoggyClubApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = _createRouter(ref);
    
    return MaterialApp.router(
      title: 'DoggyClub',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        useMaterial3: true,
      ),
      routerConfig: router,
    );
  }

  GoRouter _createRouter(WidgetRef ref) {
    return GoRouter(
      initialLocation: '/login',
      redirect: (context, state) {
        final isLoggedIn = ref.read(isLoggedInProvider);
        final isLoading = ref.read(isLoadingProvider);
        
        // Show loading while checking auth status
        if (isLoading) return null;
        
        // Redirect to home if logged in and trying to access auth pages
        if (isLoggedIn && (state.matchedLocation == '/login' || state.matchedLocation == '/register')) {
          return '/home';
        }
        
        // Redirect to login if not logged in and trying to access protected pages
        if (!isLoggedIn && (state.matchedLocation.startsWith('/home') || 
                           state.matchedLocation.startsWith('/profile') || 
                           state.matchedLocation.startsWith('/add-dog'))) {
          return '/login';
        }
        
        return null;
      },
      routes: [
        GoRoute(
          path: '/login',
          builder: (context, state) => const LoginScreen(),
        ),
        GoRoute(
          path: '/register',
          builder: (context, state) => const RegisterScreen(),
        ),
        GoRoute(
          path: '/home',
          builder: (context, state) => const HomeScreen(),
        ),
        GoRoute(
          path: '/profile',
          builder: (context, state) => const ProfileScreen(),
        ),
        GoRoute(
          path: '/add-dog',
          builder: (context, state) => const AddDogScreen(),
        ),
      ],
    );
  }
}

// Temporary home screen
class HomeScreen extends ConsumerWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(userProvider);
    
    return Scaffold(
      appBar: AppBar(
        title: const Text('DoggyClub'),
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () async {
              await ref.read(authProvider.notifier).logout();
              if (context.mounted) {
                context.go('/login');
              }
            },
          ),
        ],
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              'Welcome, ${user?.username ?? 'User'}!',
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            const SizedBox(height: 16),
            const Text('You are now logged in to DoggyClub'),
            const SizedBox(height: 32),
            ElevatedButton(
              onPressed: () {
                context.push('/profile');
              },
              child: const Text('Go to Profile'),
            ),
          ],
        ),
      ),
    );
  }
}