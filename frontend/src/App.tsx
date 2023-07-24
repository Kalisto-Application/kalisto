import MainLayout from "./layout";
import { Route, Routes, MemoryRouter } from "react-router-dom";
import { ApiPage } from "./pages/ApiPage";
import { VariablesPage } from "./pages/VariablesPage";
import {ContextProvider} from './state/ContextProvider';


function App() {
  return (
    <MemoryRouter initialEntries={['/api']} initialIndex={0}>
      <ContextProvider>
        <Routes>
          <Route path="/" element={<MainLayout />}>
            <Route path='/api' element={<ApiPage />} />
            <Route path='/variables' element={<VariablesPage />} />
          </Route>
        </Routes>
      </ContextProvider>
    </MemoryRouter>
  );
}

export default App;
  