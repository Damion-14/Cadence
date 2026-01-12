import * as React from "react";
import type { RouteObject } from "react-router"

/**
 * Auto-generate React Router routes from src/pages/**
 */

type Mod = {
  default?: any;
  loader?: any;
  action?: any;
  ErrorBoundary?: any;
};

type Node = {
  seg?: string;
  page?: Mod;
  layout?: Mod;
  children: Record<string, Node>;
};

const modules = import.meta.glob("./app/**/*.{tsx,ts,jsx,js}", { eager: true }) as Record<string, Mod>;

function toParts(filepath: string): { parts: string[], isLayout: boolean, isPage: boolean } | null {
  let p = filepath.replace(/^\.\/app\//, "").replace(/\.(t|j)sx?$/, "");

  const isLayout = p === "layout" || p.endsWith("/layout") || p.endsWith("/_layout");

  const isPage = p === "page" || p.endsWith("/page") || p === "index" || p.endsWith("/index");

  if (!isLayout && !isPage) {
    return { parts: p.split("/"), isLayout: false, isPage: true };
  }

  if (isLayout) {
    if (p === "layout") return { parts: [], isLayout: true, isPage: false };
    p = p.replace(/\/(layout|_layout)$/, "");
    return { parts: p === "" ? [] : p.split("/"), isLayout: true, isPage: false };
  }

  if (isPage) {
    if (p === "page" || p === "index") return { parts: [], isLayout: false, isPage: true };
    p = p.replace(/\/(page|index)$/, "");
    return { parts: p === "" ? [] : p.split("/"), isLayout: false, isPage: true };
  }

  return null;
}

function segToPath(seg: string): string | undefined {
  if (seg === "index" || seg === "page") return undefined;
  if (/^\[\.{3}.+\]$/.test(seg)) return "*";
  if (/^\[\[.+\]\]$/.test(seg)) return `:${seg.slice(2, -2)}?`;
  if (/^\[.+\]$/.test(seg)) return `:${seg.slice(1, -1)}`;
  return seg;
}

const root: Node = { children: {} };

for (const [file, mod] of Object.entries(modules)) {
  const result = toParts(file);
  if (!result) continue;

  const { parts, isLayout, isPage } = result;

  let cur = root;
  for (const part of parts) {
    cur.children[part] ??= { seg: part, children: {} };
    cur = cur.children[part];
  }

  if (isLayout) {
    cur.layout = mod;
  } else if (isPage) {
    cur.page = mod;
  }
}

function nodeToRoutes(node: Node): RouteObject[] {
  const routes: RouteObject[] = [];

  if (node.page) {
    routes.push({
      index: true,
      element: node.page.default ? React.createElement(node.page.default) : undefined,
      loader: node.page.loader,
      action: node.page.action,
      errorElement: node.page.ErrorBoundary
        ? React.createElement(node.page.ErrorBoundary)
        : undefined,
    });
  }

  for (const [key, childNode] of Object.entries(node.children)) {
    const path = segToPath(key);
    const childRoutes = nodeToRoutes(childNode);
    if (childRoutes.length === 0) continue;

    if (childNode.layout) {
      routes.push({
        path,
        element: childNode.layout.default
          ? React.createElement(childNode.layout.default)
          : undefined,
        loader: childNode.layout.loader,
        action: childNode.layout.action,
        errorElement: childNode.layout.ErrorBoundary
          ? React.createElement(childNode.layout.ErrorBoundary)
          : undefined,
        children: childRoutes,
      });
    } else {
      routes.push({ path, children: childRoutes });
    }
  }

  return routes;
}

function generateRoutes(): RouteObject[] {
  const rootLayout = root.layout;

  if (!rootLayout) {
    throw new Error("Root layout not found. Please create src/app/layout.tsx");
  }

  const childRoutes = nodeToRoutes(root);

  childRoutes.push({
    path: "*",
    element: React.createElement("h1", null, "Not Found")
  });

  return [
    {
      path: "/",
      element: rootLayout.default
        ? React.createElement(rootLayout.default)
        : undefined,
      loader: rootLayout.loader,
      action: rootLayout.action,
      errorElement: rootLayout.ErrorBoundary
        ? React.createElement(rootLayout.ErrorBoundary)
        : undefined,
      children: childRoutes,
    },
  ];
}

export const generatedRoutes: RouteObject[] = generateRoutes();