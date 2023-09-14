import React, { useRef, useState, useEffect } from 'react';
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api';

type props = {
  value: string;
  onChange?: (value: string) => void;
  readonly?: boolean;
};

export const Editor: React.FC<props> = ({ value, onChange, readonly }) => {
  const [editor, setEditor] =
    useState<monaco.editor.IStandaloneCodeEditor | null>(null);
  const [sub, setSub] = useState<monaco.IDisposable | undefined>();
  const monacoEl = useRef(null);

  const clean = () => {
    sub?.dispose();
    editor?.dispose();
  };

  useEffect(() => {
    if (monacoEl) {
      setEditor((editor) => {
        if (editor && (!editor.getModel()?.isDisposed() || true)) return editor;

        monaco.languages.typescript.javascriptDefaults.setDiagnosticsOptions({
          // noSemanticValidation: true,
          noSyntaxValidation: true,
        });

        const ed = monaco.editor.create(monacoEl.current!, {
          value: value,
          language: 'javascript',
          theme: 'vs-dark',
          minimap: { enabled: false },
          readOnly: readonly,
          scrollBeyondLastLine: false,
        });
        setSub(
          ed.getModel()?.onDidChangeContent((e) => {
            if (onChange) {
              onChange(ed.getModel()?.getValue() || '');
            }
          }),
        );

        return ed;
      });
    }

    return clean;
  }, [monacoEl.current]);

  return <div className="h-[600px] w-full" ref={monacoEl}></div>;
};
