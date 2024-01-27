import { AsyncLocalStorage } from "async_hooks";

const context = new AsyncLocalStorage<Map<any, any>>();

export default context;
