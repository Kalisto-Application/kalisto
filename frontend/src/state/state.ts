import { Dispatch, createContext } from "react";
import {models} from '../../wailsjs/go/models';

type Action = 
    | {type: 'switchEditor', i: number}
    | {type: 'changeRequestText', text: string}
    | {type: 'changeMetaText', text: string}
    | {type: 'newWorkspace', workspace: models.Workspace}
    | {type: 'activeWorkspace', workspace: models.Workspace}
    | {type: 'removeWorkspace', id: string}
    | {type: 'workspaceList', workspaceList: models.Workspace[]}
    | {type: 'activeMethod', activeMethod: models.Method}
    | {type: 'changeVariables', text: string}

export type State = {
    activeEditor: number;
    requestText: string;
    requestMetaText: string;
    workspaceList: models.Workspace[];
    activeWorkspace?: models.Workspace;
    activeMethod?: models.Method;
    vars: string;
}

export const reducer = (state: State, action: Action): State => {
    switch (action.type) {
        case 'switchEditor':
            return {
                ...state,
                activeEditor: action.i as number,
            }
        case 'changeRequestText':
            return {
                ... state,
                requestText: action.text,
            }
        case 'changeMetaText':
            return {
                ... state,
                requestMetaText: action.text,
            }
        case 'newWorkspace':
            return {
                ... state,
                workspaceList: state.workspaceList.concat([action.workspace]),
                activeWorkspace: action.workspace,
            } 
        case 'activeWorkspace':
            return {
                ... state,
                activeWorkspace: action.workspace,
            }
        case 'removeWorkspace':
            return {
                ... state,
                workspaceList: state.workspaceList.filter(it => it.id != action.id),
                activeWorkspace: action.id === state.activeWorkspace?.id ? undefined: state.activeWorkspace,
            }
        case 'workspaceList':
            return {
                ... state,
                workspaceList: action.workspaceList,
                activeWorkspace: action.workspaceList.length > 0 ? action.workspaceList[0]: undefined
            }
        case 'activeMethod':
            return {
                ... state,
                activeMethod: action.activeMethod,
                requestText: action.activeMethod.requestExample,
                activeEditor: 0,
            }
        case 'changeVariables':
            return {
                ... state,
                vars: action.text,
            }
        default:
            return state
    }
}

type AppContext = {
    state: State,
    dispatch: Dispatch<Action>
}

export const Context = createContext<AppContext>({} as AppContext);
