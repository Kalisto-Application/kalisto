import { Popover, Transition } from '@headlessui/react';

type popoverUiType = {
  children: React.ReactNode;
  text: string;
};

const PopoverUI: React.FC<popoverUiType> = ({ children, text }) => {
  return (
    <Popover className="relative">
      {({ close }) => (
        <>
          <Popover.Button
            as="div"
            onClick={() => setTimeout(() => close(), 1200)}
            className="flex"
          >
            {children}
          </Popover.Button>

          <Transition
            enter="transition duration-150 ease-out"
            enterFrom="transform scale-95 opacity-0"
            enterTo="transform scale-100 opacity-100"
            leave="transition duration-300 ease-out"
            leaveFrom="transform scale-100 opacity-100"
            leaveTo="transform scale-95 opacity-0"
          >
            <Popover.Panel className="absolute -left-4">
              <div className="rounded bg-primaryFill p-1">{text}</div>
            </Popover.Panel>
          </Transition>
        </>
      )}
    </Popover>
  );
};

export default PopoverUI;
