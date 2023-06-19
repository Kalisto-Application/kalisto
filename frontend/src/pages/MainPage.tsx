import React from "react";
import MainLayout from "../layout/MainLayout";
import { UrlInput } from "../components/UrlInput";

// TODO: use cases
// import { SpecFromProto, SendGrpc } from "../../wailsjs/go/api/Api";

export const MainPage: React.FC = () => {
  return (
    <MainLayout>
      <Content />
    </MainLayout>
  );
};

const Content = () => {
  const sendRequest = (event: React.SyntheticEvent, data: string) => {
    alert(data);
  };

  return (
    <div className="p-4">
      <UrlInput onClick={sendRequest} />
    </div>
  );
};
