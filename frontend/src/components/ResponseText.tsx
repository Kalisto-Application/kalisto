import React, { useContext } from 'react';
import { Context } from '../state';
import { CodeViewer } from '../ui/Editor';
import { EditorSwitcher } from '../ui/EditorSwitcher';

export const ResponseText: React.FC = () => {
  const ctx = useContext(Context);

  const switchResponseEditor = (i: number) =>
    void [ctx.dispatch({ type: 'switchResponseEditor', i: i })];
  const bodyKey = `${ctx.state.activeMethod?.fullName}:${ctx.state.response?.body}`;
  const metaKey = `${ctx.state.activeMethod?.fullName}:${ctx.state.response?.metaData}`;

  const numberEditor = ctx.state.activeResponseEditor;

  const contentEditor =
    numberEditor === 0
      ? ctx.state.response?.body
      : numberEditor === 1
      ? ctx.state.response?.metaData
      : '';

  const editors = [
    <CodeViewer
      key={bodyKey}
      fileId={ctx.state.response?.body || ''}
      text={ctx.state.response?.body || ''}
    />,
    <CodeViewer
      key={metaKey}
      fileId={ctx.state.response?.metaData || ''}
      text={ctx.state.response?.metaData || ''}
    />,
  ];

  return (
    <div className="w-1/2 bg-textBlockFill">
      <EditorSwitcher
        items={[
          { title: 'Response', onClick: switchResponseEditor },
          { title: 'Metadata', onClick: switchResponseEditor },
        ]}
        active={ctx.state.activeResponseEditor || 0}
        onClickCopy={() => navigator.clipboard.writeText(contentEditor || '')}
      />
      {editors[ctx.state.activeResponseEditor] || editors[0]}
    </div>
  );
};
