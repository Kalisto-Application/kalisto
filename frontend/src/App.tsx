import * as Sentry from '@sentry/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import config from './config';
import MainLayout from './layout';
import { ApiPage } from './pages/ApiPage';
import { ScriptingPage } from './pages/ScriptingPage';
import { VariablesPage } from './pages/VariablesPage';
import { ContextProvider } from './state/ContextProvider';

Sentry.init({
  dsn: config.sentryDsn,
  integrations: [
    new Sentry.BrowserTracing({
      tracePropagationTargets: ['localhost'],
    }),
    new Sentry.Replay(),
  ],
  // Performance Monitoring
  tracesSampleRate: 1.0, // Capture 100% of the transactions, reduce in production!
  // Session Replay
  replaysSessionSampleRate: 0.1, // This sets the sample rate at 10%. You may want to change it to 100% while in development and then sample at a lower rate in production.
  replaysOnErrorSampleRate: 1.0, // If you're not already sampling the entire session, change the sample rate to 100% when sampling sessions where errors occur.
});

function App() {
  return (
    <MemoryRouter initialEntries={['/api']} initialIndex={0}>
      <ContextProvider>
        <Routes>
          <Route path="/" element={<MainLayout />}>
            <Route path="/api" element={<ApiPage />} />
            <Route path="/variables" element={<VariablesPage />} />
            <Route path="/scripting" element={<ScriptingPage />} />
          </Route>
        </Routes>
      </ContextProvider>
    </MemoryRouter>
  );
}

export default App;
