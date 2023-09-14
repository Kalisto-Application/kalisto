import React, { useContext, useState, useMemo, useEffect } from "react";

import { UrlInput } from "../ui/UrlInput";
import { Context } from "../state";
import { UpdateWorkspace, SendGrpc } from "../../wailsjs/go/api/Api";
import { models } from "../../wailsjs/go/models";
import { Action, debounce } from "../pkg";

export const ApiRequestSender: React.FC = () => {
  const ctx = useContext(Context);

  const activeWorkspace = ctx.state.workspaceList?.find(
    (it) => it.id === ctx.state.activeWorkspaceId,
  );
  const [url, setUrl] = useState(activeWorkspace?.targetUrl || "");

  useEffect(() => {
    setUrl(activeWorkspace?.targetUrl || "");
  }, [activeWorkspace]);

  const action: Action = (url: string) => {
    UpdateWorkspace(
      new models.Workspace({ ...activeWorkspace, targetUrl: url }),
    ).catch((err) => {
      console.log("failed to save the workspace url: ", err);
    });
  };
  let debouncedUrlUpdate: Action = useMemo<Action>(() => {
    return debounce(action, 400);
  }, [ctx.state.activeWorkspaceId]);

  const onSetUrl = (url: string) => {
    setUrl(url);
    debouncedUrlUpdate(url);
  };

  const sendRequest = (_: React.SyntheticEvent) => {
    if (!ctx.state.activeMethod || !activeWorkspace) {
      return;
    }

    SendGrpc({
      addr: url,
      workspaceId: activeWorkspace.id,
      method: ctx.state.activeMethod.fullName,
      body: ctx.state.requestText,
      meta: ctx.state.requestMetaText,
    })
      .then((res) => {
        ctx.dispatch({ type: "apiResponse", response: res });
      })
      .catch((err) => {
        if (err?.Code == "SYNTAX_ERROR") {
          ctx.dispatch({ type: "apiError", value: err.Value });
        } else if (err?.Code == "SERVER_UNAVAILABLE") {
          ctx.dispatch({ type: "apiError", value: "Server unavailable" });
        }
        console.log("failed to get response: ", err);
      });
  };

  return (
    <React.Fragment>
      <UrlInput value={url} setValue={onSetUrl} onClick={sendRequest} />
    </React.Fragment>
  );
};
