import React, { useContext, useState, useMemo, useEffect } from 'react';

import { UrlInput } from '../ui/UrlInput';
import { Context } from '../state';
import { UpdateWorkspace, RunScript } from '../../wailsjs/go/api/Api';
import { models } from '../../wailsjs/go/models';
import { Action, debounce } from '../pkg';

export const ScriptSender: React.FC = () => {
  const ctx = useContext(Context);

  const activeWorkspace = ctx.state.activeWorkspace;
  const [url, setUrl] = useState(activeWorkspace?.targetUrl || '');

  useEffect(() => {
    setUrl(activeWorkspace?.targetUrl || '');
  }, [activeWorkspace]);

  const action: Action = (url: string) => {
    UpdateWorkspace(
      new models.Workspace({ ...activeWorkspace, targetUrl: url }),
    ).catch((err) => {
      console.log('failed to save the workspace url: ', err);
    });
  };
  let debouncedUrlUpdate = useMemo<Action>(() => {
    return debounce(action, 400);
  }, [activeWorkspace]);

  const onSetUrl = (url: string) => {
    setUrl(url);
    debouncedUrlUpdate(url);
  };

  const sendRequest = (_: React.SyntheticEvent) => {
    if (!ctx.state.activeMethod || !activeWorkspace) {
      return;
    }

    RunScript({
      addr: url,
      workspaceId: activeWorkspace.id,
      body: ctx.state.scriptText,
      meta: '',
    })
      .then((res) => {
        ctx.dispatch({ type: 'scriptResponse', response: res });
      })
      .catch((err) => {
        if (err?.Code == 'SYNTAX_ERROR') {
          ctx.dispatch({ type: 'scriptError', value: err.Value });
        }
        console.log('failed to get response: ', err);
      });
  };

  return (
    <React.Fragment>
      <UrlInput value={url} setValue={onSetUrl} onClick={sendRequest} />
    </React.Fragment>
  );
};
