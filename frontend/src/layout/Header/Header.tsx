import React from "react";

type HeaderProps = {
  children?: React.ReactNode;
};

export const Header: React.FC<HeaderProps> = ({ children }) => (
  <header className="h-[92px] select-none border-b border-solid border-layoutBorder bg-transparent">
    {children}
  </header>
);
