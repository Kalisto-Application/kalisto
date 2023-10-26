import React from 'react';

import kalistoIcon from '../../../assets/icons/kalisto.svg';

export const Header: React.FC = () => {
  return (
    <header className="h-[68px] border-b-[1px] border-b-borderFill">
      <div className="h-[68px]">
        <img className="p-4" src={kalistoIcon}></img>
      </div>
    </header>
  );
};
