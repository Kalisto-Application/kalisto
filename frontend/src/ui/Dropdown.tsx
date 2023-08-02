import React, { Fragment, useState, useRef } from 'react';
import {Menu as M, MenuItem} from './Menu';
import { Menu } from '@headlessui/react';

export type DropdownProps = {
  main: DropdownItem;
  items: DropdownItem[];
}

export type DropdownItem = {
  text: string,
  icon?: string,
  tip?: tip,
  onClick?: () => void,
  edit?: editable,
  menu?: DropdownItem[],
};

interface editable {
  inEdit: boolean;
  onEditDone: (newName: string) => void;
}

type tip = 'string' | React.ReactNode;

const Dropdown: React.FC<DropdownProps> = ({main, items}) =>  {
  const [editing, setEditing] = useState("");

  return  (
    <Menu>
      <Menu.Button className="w-full p-1.5 border-[1px] border-borderFill">{main.text}</Menu.Button>
      <div className="flex relative">
      <Menu.Items className="w-full flex flex-col content-normal border-[1px] border-borderFill absolute z-10 bg-primaryFill cursor-pointer">
      {items.map((it, i) => (
          <Menu.Item key={i} as={Fragment}> 
              {(render) => {
                if (it.edit?.inEdit) {
                  setEditing(it.text)
                  return <div>
                    <input value={editing} onChange={e => setEditing(e.target.value)} autoFocus={true} onKeyDown={e => it.edit?.onEditDone(editing)} />
                  </div>
                }
                return <div className="pl-[32px] m-[4px] h-[24px] w-full ui-active:bg-textBlockFill ui-active:text-white relative">
                  <span onClick={it.onClick}>
                    {it.icon && <img src={it.icon} />}
                    <a>{it.text}</a>
                    {it.tip && typeof it.tip === 'string' ? <img src={it.tip} /> : it.tip }
                  </span>
                  {render.active && <div className='absolute z-20 left-[200px] top-0'>
                   {it.menu && <M items={it.menu.map(menu => ({text: menu.text, icon: menu.icon, onClick: menu.onClick}))}/>}
                  </div>}
                </div>}
              }
          </Menu.Item>
        ))}
      </Menu.Items>
      </div>
    </Menu>
  )
}

export default Dropdown;