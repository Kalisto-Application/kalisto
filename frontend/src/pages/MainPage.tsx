import React from "react";
import MainLayout from "../layout/MainLayout";

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
  function greet() {
    alert("hey dude");
  }

  return (
    <div className="p-8">
      <button
        className="rounded bg-buttonPrimary px-4 py-2 font-bold"
        onClick={greet}
      >
        Hey Dude
      </button>
    </div>
  );
};
