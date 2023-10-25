import React, { useContext } from 'react';
import { CodeEditor } from '../ui/Editor';
import { EditorSwitcher } from '../ui/EditorSwitcher';

import { Context } from '../state';

interface RequestEditorProps {}

export const RequestEditor: React.FC<RequestEditorProps> = () => {
  const ctx = useContext(Context);
  const switchRequestEditor = (i: number) =>
    void [ctx.dispatch({ type: 'switchRequestEditor', i: i })];

  const onChangeRequest = (text: string) => {
    ctx.dispatch({ type: 'changeRequestText', text: text });
  };
  const onChangeMeta = (text: string) => {
    ctx.dispatch({ type: 'changeMetaText', text: text });
  };

  const editors = [
    <CodeEditor
      key={0}
      text={ctx.state.requestText}
      fileId={ctx.state.activeMethod?.fullName || ''}
      onChange={onChangeRequest}
    />,
    <CodeEditor
      key={1}
      text={ctx.state.requestMetaText}
      fileId={ctx.state.activeMethod?.fullName || ''}
      onChange={onChangeMeta}
    />,
  ];

  return (
    <div className="w-1/2 bg-textBlockFill">
      <EditorSwitcher
        items={[
          { title: 'Request', onClick: switchRequestEditor },
          { title: 'Metadata', onClick: switchRequestEditor },
        ]}
        active={ctx.state.activeRequestEditor || 0}
      />
      {editors[ctx.state.activeRequestEditor] || editors[0]}
    </div>
  );
};
