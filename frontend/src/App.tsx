import { useReducer, createContext } from "react";
import { MainPage } from "./pages/MainPage";
import {reducer, State, Context} from './state';

function App() {
  const [state, dispatch] = useReducer(reducer, {} as State)

  return (
    <Context.Provider value={{state: state, dispatch: dispatch}}>
      <MainPage />
    </Context.Provider>
  );
}

export default App;
  