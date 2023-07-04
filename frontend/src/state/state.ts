import { Dispatch, createContext, useReducer } from "react";

type Action = 
    | {type: 'switchEditor', i: number}
    | {type: 'changeRequestText', text: string}
    | {type: 'changeMetaText', text: string}

export type State = {
    activeEditor: number;
    requestText: string;
    requestMetaText: string;
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
        default:
            return state
    }
}

type AppContext = {
    state: State,
    dispatch: Dispatch<Action>
}

export const Context = createContext<AppContext>({} as AppContext);
