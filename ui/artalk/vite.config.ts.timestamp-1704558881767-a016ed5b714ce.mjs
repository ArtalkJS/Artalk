// vite.config.ts
import { defineConfig } from "file:///Users/rin/WorkSpace/Artalk/node_modules/.pnpm/vite@5.0.10_@types+node@20.10.5_sass@1.69.5_terser@5.26.0/node_modules/vite/dist/node/index.js";
import tsconfigPaths from "file:///Users/rin/WorkSpace/Artalk/node_modules/.pnpm/vite-tsconfig-paths@4.2.2_typescript@5.3.3_vite@5.0.10/node_modules/vite-tsconfig-paths/dist/index.mjs";
import { resolve, dirname } from "node:path";
import { fileURLToPath } from "node:url";
import checker from "file:///Users/rin/WorkSpace/Artalk/node_modules/.pnpm/vite-plugin-checker@0.6.2_eslint@8.56.0_stylelint@16.1.0_typescript@5.3.3_vite@5.0.10/node_modules/vite-plugin-checker/dist/esm/main.js";
import dts from "file:///Users/rin/WorkSpace/Artalk/node_modules/.pnpm/vite-plugin-dts@3.6.4_@types+node@20.10.5_rollup@4.9.1_typescript@5.3.3_vite@5.0.10/node_modules/vite-plugin-dts/dist/index.mjs";
import { copyFileSync } from "node:fs";
var __vite_injected_original_import_meta_url = "file:///Users/rin/WorkSpace/Artalk/ui/artalk/vite.config.ts";
var __dirname = dirname(fileURLToPath(__vite_injected_original_import_meta_url));
function getFileName(name2, format) {
  if (format == "umd")
    return `${name2}.js`;
  else if (format == "cjs")
    return `${name2}.cjs`;
  else if (format == "es")
    return `${name2}.mjs`;
  return `${name2}.${format}.js`;
}
var name = process.env.ARTALK_LITE ? "ArtalkLite" : "Artalk";
var vite_config_default = defineConfig({
  root: __dirname,
  build: {
    target: "es2015",
    outDir: resolve(__dirname, "dist"),
    minify: "terser",
    sourcemap: true,
    emptyOutDir: name === "Artalk",
    // wait for https://github.com/qmhc/vite-plugin-dts/pull/291
    lib: {
      name,
      fileName: (format) => getFileName(name, format),
      entry: resolve(__dirname, "src/main.ts"),
      formats: ["es", "umd", "cjs", "iife"]
    },
    rollupOptions: {
      external: name === "ArtalkLite" ? ["marked"] : [],
      output: {
        globals: name === "ArtalkLite" ? {
          marked: "marked"
        } : {},
        assetFileNames: (assetInfo) => /\.css$/.test(assetInfo.name || "") ? `${name}.css` : "[name].[ext]",
        // @see https://github.com/rollup/rollup/issues/587
        //  and https://github.com/rollup/rollup/pull/631/files
        exports: "named"
      }
    }
  },
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `@import "./src/style/_variables.scss";@import "./src/style/_extend.scss";`
      }
    }
  },
  resolve: {
    alias: {
      "@": resolve(__dirname, "src"),
      "~": resolve(__dirname)
    }
  },
  define: {
    ARTALK_LITE: false
  },
  plugins: [
    tsconfigPaths(),
    checker({
      typescript: true,
      eslint: {
        lintCommand: 'eslint "./src/**/*.{js,ts}"'
      }
    }),
    // @see https://github.com/qmhc/vite-plugin-dts
    name === "Artalk" ? dts({
      include: ["src"],
      exclude: ["src/**/*.{spec,test}.ts", "dist"],
      rollupTypes: true,
      afterBuild: () => {
        copyFileSync("dist/main.d.ts", "dist/main.d.cts");
      }
    }) : null
  ]
});
export {
  vite_config_default as default,
  getFileName
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCIvVXNlcnMvcmluL1dvcmtTcGFjZS9BcnRhbGsvdWkvYXJ0YWxrXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ZpbGVuYW1lID0gXCIvVXNlcnMvcmluL1dvcmtTcGFjZS9BcnRhbGsvdWkvYXJ0YWxrL3ZpdGUuY29uZmlnLnRzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy9Vc2Vycy9yaW4vV29ya1NwYWNlL0FydGFsay91aS9hcnRhbGsvdml0ZS5jb25maWcudHNcIjtpbXBvcnQgeyBkZWZpbmVDb25maWcgfSBmcm9tICd2aXRlJ1xuaW1wb3J0IHRzY29uZmlnUGF0aHMgZnJvbSAndml0ZS10c2NvbmZpZy1wYXRocydcbmltcG9ydCB7IHJlc29sdmUsIGRpcm5hbWUgfSBmcm9tICdub2RlOnBhdGgnXG5pbXBvcnQgeyBmaWxlVVJMVG9QYXRoIH0gZnJvbSAnbm9kZTp1cmwnXG5pbXBvcnQgY2hlY2tlciBmcm9tICd2aXRlLXBsdWdpbi1jaGVja2VyJ1xuaW1wb3J0IGR0cyBmcm9tICd2aXRlLXBsdWdpbi1kdHMnXG5pbXBvcnQgeyBjb3B5RmlsZVN5bmMgfSBmcm9tIFwibm9kZTpmc1wiXG5cbmNvbnN0IF9fZGlybmFtZSA9IGRpcm5hbWUoZmlsZVVSTFRvUGF0aChpbXBvcnQubWV0YS51cmwpKVxuXG5leHBvcnQgZnVuY3Rpb24gZ2V0RmlsZU5hbWUobmFtZTogc3RyaW5nLCBmb3JtYXQ6IHN0cmluZykge1xuICBpZiAoZm9ybWF0ID09IFwidW1kXCIpIHJldHVybiBgJHtuYW1lfS5qc2BcbiAgZWxzZSBpZiAoZm9ybWF0ID09IFwiY2pzXCIpIHJldHVybiBgJHtuYW1lfS5janNgXG4gIGVsc2UgaWYgKGZvcm1hdCA9PSBcImVzXCIpIHJldHVybiBgJHtuYW1lfS5tanNgXG4gIHJldHVybiBgJHtuYW1lfS4ke2Zvcm1hdH0uanNgXG59XG5cbmNvbnN0IG5hbWUgPSBwcm9jZXNzLmVudi5BUlRBTEtfTElURSA/ICdBcnRhbGtMaXRlJyA6ICdBcnRhbGsnXG5cbmV4cG9ydCBkZWZhdWx0IGRlZmluZUNvbmZpZyh7XG4gIHJvb3Q6IF9fZGlybmFtZSxcbiAgYnVpbGQ6IHtcbiAgICB0YXJnZXQ6ICdlczIwMTUnLFxuICAgIG91dERpcjogcmVzb2x2ZShfX2Rpcm5hbWUsIFwiZGlzdFwiKSxcbiAgICBtaW5pZnk6ICd0ZXJzZXInLFxuICAgIHNvdXJjZW1hcDogdHJ1ZSxcbiAgICBlbXB0eU91dERpcjogbmFtZSA9PT0gJ0FydGFsaycsICAvLyB3YWl0IGZvciBodHRwczovL2dpdGh1Yi5jb20vcW1oYy92aXRlLXBsdWdpbi1kdHMvcHVsbC8yOTFcbiAgICBsaWI6IHtcbiAgICAgIG5hbWU6IG5hbWUsXG4gICAgICBmaWxlTmFtZTogKGZvcm1hdDogc3RyaW5nKSA9PiBnZXRGaWxlTmFtZShuYW1lLCBmb3JtYXQpLFxuICAgICAgZW50cnk6IHJlc29sdmUoX19kaXJuYW1lLCAnc3JjL21haW4udHMnKSxcbiAgICAgIGZvcm1hdHM6IFtcImVzXCIsIFwidW1kXCIsIFwiY2pzXCIsIFwiaWlmZVwiXVxuICAgIH0sXG4gICAgcm9sbHVwT3B0aW9uczoge1xuICAgICAgZXh0ZXJuYWw6IChuYW1lID09PSAnQXJ0YWxrTGl0ZScpID8gWydtYXJrZWQnXSA6IFtdLFxuICAgICAgb3V0cHV0OiB7XG4gICAgICAgIGdsb2JhbHM6IChuYW1lID09PSAnQXJ0YWxrTGl0ZScpID8ge1xuICAgICAgICAgIG1hcmtlZDogJ21hcmtlZCcsXG4gICAgICAgIH0gOiB7fSxcbiAgICAgICAgYXNzZXRGaWxlTmFtZXM6IChhc3NldEluZm8pID0+ICgvXFwuY3NzJC8udGVzdChhc3NldEluZm8ubmFtZSB8fCAnJykgPyBgJHtuYW1lfS5jc3NgIDogXCJbbmFtZV0uW2V4dF1cIiksXG4gICAgICAgIC8vIEBzZWUgaHR0cHM6Ly9naXRodWIuY29tL3JvbGx1cC9yb2xsdXAvaXNzdWVzLzU4N1xuICAgICAgICAvLyAgYW5kIGh0dHBzOi8vZ2l0aHViLmNvbS9yb2xsdXAvcm9sbHVwL3B1bGwvNjMxL2ZpbGVzXG4gICAgICAgIGV4cG9ydHM6ICduYW1lZCcsXG4gICAgICB9XG4gICAgfVxuICB9LFxuICBjc3M6IHtcbiAgICBwcmVwcm9jZXNzb3JPcHRpb25zOiB7XG4gICAgICBzY3NzOiB7XG4gICAgICAgIGFkZGl0aW9uYWxEYXRhOiBgQGltcG9ydCBcIi4vc3JjL3N0eWxlL192YXJpYWJsZXMuc2Nzc1wiO0BpbXBvcnQgXCIuL3NyYy9zdHlsZS9fZXh0ZW5kLnNjc3NcIjtgXG4gICAgIH0sXG4gICAgfSxcbiAgfSxcbiAgcmVzb2x2ZToge1xuICAgIGFsaWFzOiB7XG4gICAgICAnQCc6IHJlc29sdmUoX19kaXJuYW1lLCAnc3JjJyksXG4gICAgICAnfic6IHJlc29sdmUoX19kaXJuYW1lKSxcbiAgICB9XG4gIH0sXG4gIGRlZmluZToge1xuICAgIEFSVEFMS19MSVRFOiBmYWxzZSxcbiAgfSxcbiAgcGx1Z2luczogW1xuICAgIHRzY29uZmlnUGF0aHMoKSxcbiAgICBjaGVja2VyKHtcbiAgICAgIHR5cGVzY3JpcHQ6IHRydWUsXG4gICAgICBlc2xpbnQ6IHtcbiAgICAgICAgbGludENvbW1hbmQ6ICdlc2xpbnQgXCIuL3NyYy8qKi8qLntqcyx0c31cIicsXG4gICAgICB9LFxuICAgIH0pLFxuICAgIC8vIEBzZWUgaHR0cHM6Ly9naXRodWIuY29tL3FtaGMvdml0ZS1wbHVnaW4tZHRzXG4gICAgKG5hbWUgPT09ICdBcnRhbGsnKSA/IGR0cyh7XG4gICAgICBpbmNsdWRlOiBbJ3NyYyddLFxuICAgICAgZXhjbHVkZTogWydzcmMvKiovKi57c3BlYyx0ZXN0fS50cycsICdkaXN0J10sXG4gICAgICByb2xsdXBUeXBlczogdHJ1ZSxcbiAgICAgIGFmdGVyQnVpbGQ6ICgpID0+IHtcbiAgICAgICAgLy8gQHNlZSBodHRwczovL2dpdGh1Yi5jb20vYXJldGhldHlwZXN3cm9uZy9hcmV0aGV0eXBlc3dyb25nLmdpdGh1Yi5pby90cmVlL21haW4vcGFja2FnZXMvY2xpXG4gICAgICAgIC8vIEBmaXggaHR0cHM6Ly9naXRodWIuY29tL2FyZXRoZXR5cGVzd3JvbmcvYXJldGhldHlwZXN3cm9uZy5naXRodWIuaW8vYmxvYi9tYWluL2RvY3MvcHJvYmxlbXMvRmFsc2VFU00ubWQjY29uc2VxdWVuY2VzXG4gICAgICAgIGNvcHlGaWxlU3luYyhcImRpc3QvbWFpbi5kLnRzXCIsIFwiZGlzdC9tYWluLmQuY3RzXCIpXG4gICAgICB9LFxuICAgIH0pIDogbnVsbCxcbiAgXSxcbn0pXG4iXSwKICAibWFwcGluZ3MiOiAiO0FBQWlTLFNBQVMsb0JBQW9CO0FBQzlULE9BQU8sbUJBQW1CO0FBQzFCLFNBQVMsU0FBUyxlQUFlO0FBQ2pDLFNBQVMscUJBQXFCO0FBQzlCLE9BQU8sYUFBYTtBQUNwQixPQUFPLFNBQVM7QUFDaEIsU0FBUyxvQkFBb0I7QUFOcUosSUFBTSwyQ0FBMkM7QUFRbk8sSUFBTSxZQUFZLFFBQVEsY0FBYyx3Q0FBZSxDQUFDO0FBRWpELFNBQVMsWUFBWUEsT0FBYyxRQUFnQjtBQUN4RCxNQUFJLFVBQVU7QUFBTyxXQUFPLEdBQUdBLEtBQUk7QUFBQSxXQUMxQixVQUFVO0FBQU8sV0FBTyxHQUFHQSxLQUFJO0FBQUEsV0FDL0IsVUFBVTtBQUFNLFdBQU8sR0FBR0EsS0FBSTtBQUN2QyxTQUFPLEdBQUdBLEtBQUksSUFBSSxNQUFNO0FBQzFCO0FBRUEsSUFBTSxPQUFPLFFBQVEsSUFBSSxjQUFjLGVBQWU7QUFFdEQsSUFBTyxzQkFBUSxhQUFhO0FBQUEsRUFDMUIsTUFBTTtBQUFBLEVBQ04sT0FBTztBQUFBLElBQ0wsUUFBUTtBQUFBLElBQ1IsUUFBUSxRQUFRLFdBQVcsTUFBTTtBQUFBLElBQ2pDLFFBQVE7QUFBQSxJQUNSLFdBQVc7QUFBQSxJQUNYLGFBQWEsU0FBUztBQUFBO0FBQUEsSUFDdEIsS0FBSztBQUFBLE1BQ0g7QUFBQSxNQUNBLFVBQVUsQ0FBQyxXQUFtQixZQUFZLE1BQU0sTUFBTTtBQUFBLE1BQ3RELE9BQU8sUUFBUSxXQUFXLGFBQWE7QUFBQSxNQUN2QyxTQUFTLENBQUMsTUFBTSxPQUFPLE9BQU8sTUFBTTtBQUFBLElBQ3RDO0FBQUEsSUFDQSxlQUFlO0FBQUEsTUFDYixVQUFXLFNBQVMsZUFBZ0IsQ0FBQyxRQUFRLElBQUksQ0FBQztBQUFBLE1BQ2xELFFBQVE7QUFBQSxRQUNOLFNBQVUsU0FBUyxlQUFnQjtBQUFBLFVBQ2pDLFFBQVE7QUFBQSxRQUNWLElBQUksQ0FBQztBQUFBLFFBQ0wsZ0JBQWdCLENBQUMsY0FBZSxTQUFTLEtBQUssVUFBVSxRQUFRLEVBQUUsSUFBSSxHQUFHLElBQUksU0FBUztBQUFBO0FBQUE7QUFBQSxRQUd0RixTQUFTO0FBQUEsTUFDWDtBQUFBLElBQ0Y7QUFBQSxFQUNGO0FBQUEsRUFDQSxLQUFLO0FBQUEsSUFDSCxxQkFBcUI7QUFBQSxNQUNuQixNQUFNO0FBQUEsUUFDSixnQkFBZ0I7QUFBQSxNQUNuQjtBQUFBLElBQ0Q7QUFBQSxFQUNGO0FBQUEsRUFDQSxTQUFTO0FBQUEsSUFDUCxPQUFPO0FBQUEsTUFDTCxLQUFLLFFBQVEsV0FBVyxLQUFLO0FBQUEsTUFDN0IsS0FBSyxRQUFRLFNBQVM7QUFBQSxJQUN4QjtBQUFBLEVBQ0Y7QUFBQSxFQUNBLFFBQVE7QUFBQSxJQUNOLGFBQWE7QUFBQSxFQUNmO0FBQUEsRUFDQSxTQUFTO0FBQUEsSUFDUCxjQUFjO0FBQUEsSUFDZCxRQUFRO0FBQUEsTUFDTixZQUFZO0FBQUEsTUFDWixRQUFRO0FBQUEsUUFDTixhQUFhO0FBQUEsTUFDZjtBQUFBLElBQ0YsQ0FBQztBQUFBO0FBQUEsSUFFQSxTQUFTLFdBQVksSUFBSTtBQUFBLE1BQ3hCLFNBQVMsQ0FBQyxLQUFLO0FBQUEsTUFDZixTQUFTLENBQUMsMkJBQTJCLE1BQU07QUFBQSxNQUMzQyxhQUFhO0FBQUEsTUFDYixZQUFZLE1BQU07QUFHaEIscUJBQWEsa0JBQWtCLGlCQUFpQjtBQUFBLE1BQ2xEO0FBQUEsSUFDRixDQUFDLElBQUk7QUFBQSxFQUNQO0FBQ0YsQ0FBQzsiLAogICJuYW1lcyI6IFsibmFtZSJdCn0K
