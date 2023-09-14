import { ReactNode, useEffect, useReducer } from 'react';
import * as Sentry from '@sentry/react';

import { reducer, Context, newState } from './state';

import { models } from '../../wailsjs/go/models';
import { FindWorkspaces, GetGlobalVars } from '../../wailsjs/go/api/Api';

type ContextProps = {
  children?: ReactNode;
};

export const ContextProvider: React.FC<ContextProps> = ({ children }) => {
  const [state, dispatch] = useReducer(reducer, newState());

  useEffect(() => {
    // load the app
    FindWorkspaces()
      .then((res) => {
        if (res.length === 0) {
          dispatch({ type: 'workspaceList', workspaceList: [] });
          return;
        }

        let latest = res.at(0)?.lastUsage || new Date(0);
        res.forEach((it) => {
          if (it.lastUsage > latest) {
            latest = it.lastUsage;
          }
        });

        const getFirstMethod = (): models.Method | undefined => {
          for (const ws of res) {
            for (const service of ws.spec.services) {
              for (const m of service.methods) {
                return m;
              }
            }
          }
        };

        const fristMethod = getFirstMethod();
        if (fristMethod) {
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
    <Sentry.ErrorBoundary
      beforeCapture={(scope) => {
        scope.setExtra('state', state);
      }}
    >
      <Context.Provider value={{ state, dispatch }}>
        {children}
      </Context.Provider>
    </Sentry.ErrorBoundary>
  );
};
