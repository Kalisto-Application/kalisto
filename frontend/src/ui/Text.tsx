interface props {
  text: string;
  props?: any;
}

export const Text: React.FC<props> = ({ text, props }) => {
  return <span {...props}>{text}</span>;
};

export default Text;
