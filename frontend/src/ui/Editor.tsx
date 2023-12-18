import * as monaco from 'monaco-editor/esm/vs/editor/editor.api';
import React, { useEffect, useMemo, useRef, useState } from 'react';

import { debounce } from '../pkg';

type CodeEditorProps = {
  text: string;
  fileId: string;
  onChange: (text: string) => void;
};

export const CodeEditor: React.FC<CodeEditorProps> = ({
  text,
  fileId,
  onChange,
}) => {
  const handleOnChange = useMemo(() => {
    return debounce(onChange, 400);
  }, [fileId]);

  return <Editor key={fileId} value={text} onChange={handleOnChange} />;
};

type CodeViewerProps = {
  text: string;
  fileId: string;
};

export const CodeViewer: React.FC<CodeViewerProps> = ({ text, fileId }) => {
  return <Editor key={fileId} value={text} readonly />;
};

type props = {
  value: string;
  onChange?: (value: string) => void;
  readonly?: boolean;
};

const Editor: React.FC<props> = ({ value, onChange, readonly }) => {
  const [editor, setEditor] =
    useState<monaco.editor.IStandaloneCodeEditor | null>(null);
  const [sub, setSub] = useState<monaco.IDisposable | undefined>();
  const monacoEl = useRef(null);

  const clean = (): void => {
    sub?.dispose();
    editor?.dispose();
  };

  useEffect(() => {
    if (monacoEl !== undefined) {
      setEditor((editor) => {
        if (editor !== null && (!editor.getModel()?.isDisposed() || true))
          return editor;

        monaco.languages.typescript.javascriptDefaults.setDiagnosticsOptions({
          // noSemanticValidation: true,
          noSyntaxValidation: true,
        });

        const ed = monaco.editor.create(monacoEl.current!, {
          value: value,
          language: 'javascript',
          theme: 'vs-dark',
          minimap: { enabled: false },
          overviewRulerLanes: 0,
          readOnly: readonly,
          scrollBeyondLastLine: false,
        });
        setSub(
          ed.getModel()?.onDidChangeContent((e) => {
            if (onChange !== undefined) {
              onChange(ed.getModel()?.getValue() || '');
            }
          })
        );

        return ed;
      });
    }

    return clean;
  }, [monacoEl.current]);

  return (
    <div
      className="-z-1 relative  h-[77.4vh] w-full border-l-2 border-borderFill"
      ref={monacoEl}
    ></div>
  );
};
