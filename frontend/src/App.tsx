import { useReducer, useContext, useEffect } from "react";
import { FindWorkspaces } from "../wailsjs/go/api/Api";
import { models } from "../wailsjs/go/models";
import MainLayout from "./layout/MainLayout";
import {reducer, State, Context} from './state';

function App() {
  const [state, dispatch] = useReducer(reducer, {} as State);

  return (
    <Context.Provider value={{state: state, dispatch: dispatch}}>
      <MainLayout />
    </Context.Provider>
  );
}

export default App;
  