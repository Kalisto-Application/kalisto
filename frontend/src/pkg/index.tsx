export const debounce = (f: Action, delay: number): Action => {
    let timer: NodeJS.Timeout;
    return (data: any) => {
      clearTimeout(timer);
      timer = setTimeout(() => { f(data); }, delay);
    };
}

export type Action = (data: any) => void
