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
  idFile: string;
  action: Action;
}

export const CodeEditor: React.FC<CodeEditorProps> = ({
  text,
  action,
  idFile,
}) => {
  const ctx = useContext(Context);

  useEffect(() => {
    return () => {
      console.log('Компонент удален', idFile);
    };
  }, [idFile]);

  let onChange = (value: string) => {
    action(value);
  };

  return <Editor value={text} onChange={onChange} />;
};
