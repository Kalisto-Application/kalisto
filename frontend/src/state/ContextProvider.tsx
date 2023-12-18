import * as Sentry from '@sentry/react';
import { useEffect, useReducer, type ReactNode } from 'react';

import { Context, newState, reducer } from './state';

import { GetGlobalVars, WorkspaceList } from '../../wailsjs/go/api/Api';
import type { models } from '../../wailsjs/go/models';

type ContextProps = {
  children?: ReactNode;
};

export const ContextProvider: React.FC<ContextProps> = ({ children }) => {
  const [state, dispatch] = useReducer(reducer, newState());

  useEffect(() => {
    // load the app
    WorkspaceList('')
      .then((res) => {
        dispatch({ type: 'workspaceList', workspaceList: res });

        const getFirstMethod = (): models.Method | undefined => {
          for (const service of res.main.spec.services) {
            for (const m of service.methods) {
              return m;
            }
          }
        };

        const fristMethod = getFirstMethod();
        if (fristMethod !== undefined) {
          dispatch({ type: 'activeMethod', activeMethod: fristMethod });
        }
        dispatch({ type: 'workspaceList', workspaceList: res });
      })
      .catch((err) => console.log('error on find workspaces: ', err));

    GetGlobalVars()
      .then((vars) => {
        dispatch({ type: 'changeVariables', text: vars });
      })
      .catch((err) => console.log('error on get global variables: ', err));
  }, []);

  return (
    <Sentry.ErrorBoundary>
      <Context.Provider value={{ state, dispatch }}>
        {children}
      </Context.Provider>
    </Sentry.ErrorBoundary>
  );
};
