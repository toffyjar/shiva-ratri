import * as React from 'react';
import { View, StyleSheet } from 'react-native';
import { SafeAreaProvider, SafeAreaView } from 'react-native-safe-area-context';
import CreateKundli from './pages/CreateKundli';
import ViewKundli from './pages/ViewKundli';
import Login from './pages/Login';
import Register from './pages/Register';
import { ActivityIndicator, MD3LightTheme, PaperProvider } from 'react-native-paper';
import { initDatabase } from './database/database';
import { useAuthStore } from './store/authStore';
import { useSavedKundalisStore } from './store/savedKundalisStore';

// You can customize your theme using Material Design 3 (MD3) guidelines
const theme = {
  ...MD3LightTheme,
  colors: {
    ...MD3LightTheme.colors,
    primary: '#6200ee',
    secondary: '#03dac6',
  },
};

const styles = StyleSheet.create({
  base: {
    flex: 1,
  },
  centered: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
});

type AppRoute = 'create' | 'view';
type AuthRoute = 'login' | 'register';

export default function App() {
  const [route, setRoute] = React.useState<AppRoute>('create');
  const [authRoute, setAuthRoute] = React.useState<AuthRoute>('login');
  const { isAuthenticated, isHydrating, token, hydrate, logout } = useAuthStore();
  const fetchKundalis = useSavedKundalisStore((state) => state.fetchKundalis);

  React.useEffect(() => {
    initDatabase();
    hydrate();
  }, []);

  React.useEffect(() => {
    if (isAuthenticated) {
      fetchKundalis(token);
    }
  }, [isAuthenticated, token]);

  const renderContent = () => {
    if (isHydrating) {
      return (
        <View style={styles.centered}>
          <ActivityIndicator animating size="large" />
        </View>
      );
    }
    // logout(); // Ensure the user is logged out when the app starts

    if (!isAuthenticated) {
      return authRoute === 'login' ? (
        <Login
          onLoginSuccess={() => setRoute('create')}
          onNavigateToRegister={() => setAuthRoute('register')}
        />
      ) : (
        <Register
          onRegisterSuccess={() => setRoute('create')}
          onNavigateToLogin={() => setAuthRoute('login')}
        />
      );
    }

    return route === 'create' ? (
      <CreateKundli onNavigateToView={() => setRoute('view')} />
    ) : (
      <ViewKundli onNavigateToCreate={() => setRoute('create')} />
    );
  };

  return (
    <SafeAreaProvider>
      <PaperProvider theme={theme}>
        <SafeAreaView style={styles.base}>{renderContent()}</SafeAreaView>
      </PaperProvider>
    </SafeAreaProvider>
  );
}
