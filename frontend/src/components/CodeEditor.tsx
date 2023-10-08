import React, {
  ChangeEvent,
  useContext,
  useEffect,
  useMemo,
  useState,
} from 'react';
import { Context } from '../state';
import { debounce, Action } from '../pkg';

import { Editor } from '../ui/Editor';

interface CodeEditorProps {
  text: string;
  type:
    | 'changeRequestText'
    | 'changeMetaText'
    | 'changeVariables'
    | 'changeScriptText';
  action?: Action;
  showComponent?: boolean;
  idFile?: string;
}

export const CodeEditor: React.FC<CodeEditorProps> = ({
  text,
  type,
  action,
  showComponent,
  idFile,
}) => {
  const ctx = useContext(Context);

  useEffect(() => {
    return () => {
      console.log('Компонент удален', idFile);
    };
  }, [idFile]);

  let debouncedAction: Action | undefined;
  if (action) {
    debouncedAction = useMemo<Action>(() => {
      return debounce(action, 400);
    }, []);
  }

  let onChange = (value: string) => {
    ctx.dispatch({ type: type, text: value });
    if (debouncedAction) {
      debouncedAction(value);
    }
  };

  return (
    <div>
      {showComponent && (
        <Editor
          key={ctx.state.activeMethod?.fullName || ''}
          value={text}
          onChange={onChange}
        />
      )}
    </div>
  );
};
