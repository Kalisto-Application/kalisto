import { createContext, Dispatch } from 'react';
import { models } from '../../wailsjs/go/models';

type Action =
  | { type: 'switchRequestEditor'; i: number }
  | { type: 'switchResponseEditor'; i: number }
  | { type: 'changeRequestText'; text: string }
  | { type: 'changeMetaText'; text: string }
  | { type: 'newWorkspace'; workspace: models.Workspace }
  | { type: 'removeWorkspace'; id: string }
  | { type: 'renameWorkspace'; id: string; name: string }
  | { type: 'workspaceList'; workspaceList: models.WorkspaceList }
  | { type: 'activeMethod'; activeMethod: models.Method }
  | { type: 'apiResponse'; response: models.Response }
  | { type: 'apiError'; value: string }
  | { type: 'changeVariables'; text: string }
  | { type: 'varsError'; value: string }
  | { type: 'scriptResponse'; response: string }
  | { type: 'scriptError'; value: string }
  | { type: 'updateWorkspace'; workspace: models.Workspace }
  | { type: 'setActiveScriptId'; id: string }
  | { type: 'updateScriptFile'; file: models.File }
  | { type: 'addScriptFile'; file: models.File }
  | { type: 'switchScriptEditor'; i: number }
  | { type: 'addRequestFile'; file: { [key: string]: models.File[] } }
  | { type: 'updateRequestFile'; file: models.File; metName: string }
  | { type: 'setActiveRequest'; id: string };

export type State = {
  activeRequestEditor: number;
  activeResponseEditor: number;
  requestText: string;
  requestMetaText: string;
  workspaceList?: models.WorkspaceShort[];
  activeWorkspace?: models.Workspace;
  activeMethod?: models.Method;
  response?: models.Response;
  apiError: string;
  vars: string;
  varsError: string;
  scriptResponse: string;
  scriptError: string;
  activeScriptFileId: string;
  activeRequestFileId: string;
  activeScriptEditor: number;
};

export const newState = (): State => {
  return {
    activeRequestEditor: 0,
    activeResponseEditor: 0,
    requestText: '',
    requestMetaText: '',
    workspaceList: undefined,
    activeWorkspace: undefined,
    activeMethod: undefined,
    response: undefined,
    apiError: '',
    vars: '{}',
    varsError: '',
    scriptResponse: '',
    scriptError: '',
    activeScriptFileId: '',
    activeRequestFileId: '',
    activeScriptEditor: 0,
  };
};

export const reducer = (state: State, action: Action): State => {
  console.log('dispatch');
  console.log(action);

  switch (action.type) {
    case 'switchRequestEditor':
      return {
        ...state,
        activeRequestEditor: action.i as number,
      };
    case 'switchResponseEditor':
      return {
        ...state,
        activeResponseEditor: action.i as number,
      };
    case 'changeRequestText':
      return {
        ...state,
        requestText: action.text,
      };
    case 'changeMetaText':
      return {
        ...state,
        requestMetaText: action.text,
      };
    case 'newWorkspace':
      return {
        ...state,
        workspaceList: [action.workspace as models.WorkspaceShort].concat(
          state.workspaceList || []
        ),
        activeWorkspace: action.workspace,
      };
    case 'removeWorkspace':
      const filtered = state.workspaceList?.filter((it) => it.id != action.id);
      return {
        ...state,
        workspaceList: filtered,
        activeWorkspace:
          action.id === state.activeWorkspace?.id
            ? undefined
            : state.activeWorkspace,
      };
    case 'renameWorkspace':
      let ws = state.activeWorkspace;
      if (state.activeWorkspace && state.activeWorkspace?.id === action.id) {
        ws = new models.Workspace({
          ...ws,
          name: action.name,
        });
      }
      return {
        ...state,
        activeWorkspace: ws,
        workspaceList: state.workspaceList?.map((it) => {
          if (it.id === action.id) {
            it.name = action.name;
          }
          return it;
        }),
      };
    case 'workspaceList':
      return {
        ...state,
        workspaceList: action.workspaceList.list,
        activeWorkspace: action.workspaceList.main.id
          ? action.workspaceList.main
          : undefined,
      };
    case 'activeMethod':
      return {
        ...state,
        activeMethod: action.activeMethod,
        requestText: action.activeMethod?.requestExample || '',
        activeRequestEditor: 0,
      };
    case 'apiResponse':
      return {
        ...state,
        response: action.response,
        activeResponseEditor: 0,
        apiError: '',
      };
    case 'apiError':
      return {
        ...state,
        apiError: action.value,
      };
    case 'changeVariables':
      return {
        ...state,
        vars: action.text,
        varsError: '',
      };
    case 'varsError':
      return {
        ...state,
        varsError: action.value,
      };
    case 'updateScriptFile':
      return {
        ...state,
        activeWorkspace: new models.Workspace({
          ...state.activeWorkspace,
          scriptFiles: state.activeWorkspace?.scriptFiles.map((it) => {
            if (it.id !== action.file.id) return it;
            return action.file;
          }),
        }),
      };

    case 'scriptResponse':
      return {
        ...state,
        scriptResponse: action.response,
        scriptError: '',
      };
    case 'scriptError':
      return {
        ...state,
        scriptError: action.value,
      };
    case 'updateWorkspace':
      return {
        ...state,
        activeWorkspace: action.workspace,
      };
    case 'setActiveScriptId':
      return {
        ...state,
        activeScriptFileId: action.id,
        activeScriptEditor: 0,
      };
    case 'addScriptFile':
      return {
        ...state,
        activeWorkspace: new models.Workspace({
          ...state.activeWorkspace,
          scriptFiles: [
            ...(state.activeWorkspace?.scriptFiles || []),
            action.file,
          ],
        }),
      };
    case 'switchScriptEditor':
      return {
        ...state,
        activeScriptEditor: action.i,
      };
    case 'addRequestFile':
      return {
        ...state,
        activeWorkspace: new models.Workspace({
          ...state.activeWorkspace,
          requestFiles: {
            ...action.file,
          },
        }),
      };

    case 'updateRequestFile': {
      return {
        ...state,
        activeWorkspace: new models.Workspace({
          ...state.activeWorkspace,
          requestFiles: {
            ...state.activeWorkspace?.requestFiles,

            [action.metName]: state.activeWorkspace?.requestFiles[
              action.metName
            ].map((it) => {
              if (it.id === action.file.id) return action.file;
              return it;
            }),
          },
        }),
      };
    }
    case 'setActiveRequest':
      return {
        ...state,
        activeRequestFileId: action.id,
      };
    default:
      return state;
  }
};

type AppContext = {
  state: State;
  dispatch: Dispatch<Action>;
};

export const Context = createContext<AppContext>({} as AppContext);
