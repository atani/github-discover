import React from "react";
import {
  AbsoluteFill,
  useCurrentFrame,
  spring,
  useVideoConfig,
} from "remotion";

const FONT = 'Menlo, Monaco, "Courier New", monospace';

const colors = {
  bg: "#1e1e2e",
  text: "#cdd6f4",
  green: "#a6e3a1",
  yellow: "#f9e2af",
  blue: "#89b4fa",
  purple: "#cba6f7",
  gray: "#6c7086",
  border: "#45475a",
  cyan: "#94e2d5",
  red: "#f38ba8",
};

const TerminalChrome: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => (
  <div
    style={{
      width: 880,
      margin: "0 auto",
      borderRadius: 12,
      overflow: "hidden",
      boxShadow: "0 20px 60px rgba(0,0,0,0.5)",
      border: `1px solid ${colors.border}`,
    }}
  >
    <div
      style={{
        background: "#313244",
        padding: "10px 16px",
        display: "flex",
        alignItems: "center",
        gap: 8,
      }}
    >
      {[colors.red, colors.yellow, colors.green].map((c, i) => (
        <div
          key={i}
          style={{ width: 12, height: 12, borderRadius: "50%", background: c }}
        />
      ))}
      <span
        style={{
          color: colors.gray,
          fontSize: 13,
          fontFamily: FONT,
          marginLeft: 8,
        }}
      >
        github-discover
      </span>
    </div>
    <div
      style={{
        background: colors.bg,
        padding: "20px 24px",
        minHeight: 380,
        fontFamily: FONT,
        fontSize: 14,
        lineHeight: 1.6,
      }}
    >
      {children}
    </div>
  </div>
);

const TypedText: React.FC<{
  text: string;
  startFrame: number;
  color?: string;
}> = ({ text, startFrame, color = colors.text }) => {
  const frame = useCurrentFrame();
  const charsToShow = Math.min(
    Math.floor((frame - startFrame) * 1.5),
    text.length
  );
  if (frame < startFrame) return null;
  const visible = text.slice(0, Math.max(0, charsToShow));

  return (
    <span style={{ color }}>
      {visible}
      {charsToShow < text.length && (
        <span style={{ background: colors.green, color: colors.bg }}>
          {" "}
        </span>
      )}
    </span>
  );
};

const FadeInLine: React.FC<{
  children: React.ReactNode;
  startFrame: number;
}> = ({ children, startFrame }) => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();
  const opacity = spring({
    frame: frame - startFrame,
    fps,
    config: { damping: 20 },
  });
  if (frame < startFrame) return null;
  return <div style={{ opacity }}>{children}</div>;
};

const RepoRow: React.FC<{
  rank: number;
  name: string;
  stars: string;
  lang: string;
  desc: string;
  startFrame: number;
}> = ({ rank, name, stars, lang, desc, startFrame }) => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();
  const progress = spring({
    frame: frame - startFrame,
    fps,
    config: { damping: 15 },
  });
  if (frame < startFrame) return null;

  return (
    <div
      style={{
        opacity: progress,
        transform: `translateX(${(1 - progress) * 20}px)`,
      }}
    >
      <div style={{ display: "flex", gap: 8, alignItems: "baseline" }}>
        <span style={{ color: colors.gray, width: 30, textAlign: "right" }}>
          {rank}.
        </span>
        <span style={{ color: colors.green, width: 340 }}>{name}</span>
        <span style={{ color: colors.yellow }}>★ {stars}</span>
        {lang && <span style={{ color: colors.cyan }}>[{lang}]</span>}
      </div>
      <div style={{ color: colors.gray, paddingLeft: 38, fontSize: 13 }}>
        {desc}
      </div>
    </div>
  );
};

// Scene 1: trending (0-150)
const TrendingScene: React.FC = () => (
  <TerminalChrome>
    <div>
      <span style={{ color: colors.green }}>$ </span>
      <TypedText text="github-discover trending -n 5" startFrame={10} />
    </div>
    <FadeInLine startFrame={40}>
      <div style={{ marginTop: 16 }}>
        <span style={{ color: colors.purple, fontWeight: "bold" }}>
          Trending Repositories (weekly)
        </span>
      </div>
      <div style={{ color: colors.border }}>{"─".repeat(70)}</div>
    </FadeInLine>
    <RepoRow rank={1} name="freeCodeCamp/freeCodeCamp" stars="438.7k" lang="TypeScript" desc="Open-source codebase and curriculum" startFrame={50} />
    <RepoRow rank={2} name="public-apis/public-apis" stars="414.4k" lang="Python" desc="A collective list of free APIs" startFrame={58} />
    <RepoRow rank={3} name="EbookFoundation/free-programming-books" stars="384.4k" lang="Python" desc="Freely available programming books" startFrame={66} />
    <RepoRow rank={4} name="kamranahmedse/developer-roadmap" stars="351.5k" lang="TypeScript" desc="Interactive roadmaps and guides for developers" startFrame={74} />
    <RepoRow rank={5} name="donnemartin/system-design-primer" stars="339.9k" lang="Python" desc="Learn how to design large-scale systems" startFrame={82} />
    <FadeInLine startFrame={92}>
      <div style={{ color: colors.border, marginTop: 4 }}>{"─".repeat(70)}</div>
      <div style={{ color: colors.gray, fontStyle: "italic", fontSize: 13 }}>
        Use --language to filter by language, --since to change time range
      </div>
    </FadeInLine>
  </TerminalChrome>
);

// Scene 2: search ai (150-300)
const SearchScene: React.FC = () => (
  <TerminalChrome>
    <div>
      <span style={{ color: colors.green }}>$ </span>
      <TypedText text="github-discover search ai -n 5" startFrame={155} />
    </div>
    <FadeInLine startFrame={185}>
      <div style={{ marginTop: 16 }}>
        <span style={{ color: colors.purple, fontWeight: "bold" }}>
          Search Results: "ai"
        </span>
      </div>
      <div style={{ color: colors.border }}>{"─".repeat(70)}</div>
    </FadeInLine>
    <RepoRow rank={1} name="openclaw/openclaw" stars="330.5k" lang="TypeScript" desc="Your own personal AI assistant" startFrame={195} />
    <RepoRow rank={2} name="Significant-Gravitas/AutoGPT" stars="182.7k" lang="Python" desc="AutoGPT is the vision of accessible AI for everyone" startFrame={203} />
    <RepoRow rank={3} name="n8n-io/n8n" stars="180.6k" lang="TypeScript" desc="Workflow automation with native AI capabilities" startFrame={211} />
    <RepoRow rank={4} name="AUTOMATIC1111/stable-diffusion-webui" stars="161.9k" lang="Python" desc="Stable Diffusion web UI" startFrame={219} />
    <RepoRow rank={5} name="f/prompts.chat" stars="153.9k" lang="HTML" desc="Share, discover, and collect prompts" startFrame={227} />
    <FadeInLine startFrame={237}>
      <div style={{ color: colors.border, marginTop: 4 }}>{"─".repeat(70)}</div>
      <div style={{ color: colors.gray, fontStyle: "italic", fontSize: 13 }}>
        4,074,415 repositories found
      </div>
    </FadeInLine>
  </TerminalChrome>
);

// Scene 3: random (300-450)
const RandomScene: React.FC = () => (
  <TerminalChrome>
    <div>
      <span style={{ color: colors.green }}>$ </span>
      <TypedText text="github-discover random" startFrame={305} />
    </div>
    <FadeInLine startFrame={330}>
      <div
        style={{
          marginTop: 16,
          border: `1px solid ${colors.blue}`,
          borderRadius: 8,
          padding: "16px 20px",
        }}
      >
        <div style={{ color: colors.green, fontWeight: "bold", fontSize: 16 }}>
          torvalds/linux
        </div>
        <div style={{ color: colors.border, margin: "8px 0" }}>
          {"─".repeat(45)}
        </div>
        <div style={{ color: colors.text, marginBottom: 12 }}>
          Linux kernel source tree
        </div>
        <div style={{ display: "flex", flexDirection: "column", gap: 4 }}>
          <div>
            <span style={{ color: colors.cyan }}>Stars: </span>
            <span style={{ color: colors.text }}>224.6k</span>
          </div>
          <div>
            <span style={{ color: colors.cyan }}>Forks: </span>
            <span style={{ color: colors.text }}>60,142</span>
          </div>
          <div>
            <span style={{ color: colors.cyan }}>Language: </span>
            <span style={{ color: colors.text }}>C</span>
          </div>
          <div>
            <span style={{ color: colors.cyan }}>License: </span>
            <span style={{ color: colors.text }}>GPL-2.0</span>
          </div>
          <div>
            <span style={{ color: colors.cyan }}>URL: </span>
            <span style={{ color: colors.blue }}>
              https://github.com/torvalds/linux
            </span>
          </div>
        </div>
      </div>
    </FadeInLine>
    <FadeInLine startFrame={360}>
      <div
        style={{
          color: colors.gray,
          fontStyle: "italic",
          fontSize: 13,
          marginTop: 12,
        }}
      >
        Run again for different results!
      </div>
    </FadeInLine>
  </TerminalChrome>
);

export const TerminalDemo: React.FC = () => {
  const frame = useCurrentFrame();

  // Scene 1: trending (0-150), Scene 2: search (150-300), Scene 3: random (300-450)
  const scene = frame < 150 ? 1 : frame < 300 ? 2 : 3;

  return (
    <AbsoluteFill
      style={{
        background: "linear-gradient(135deg, #11111b 0%, #1e1e2e 100%)",
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      {scene === 1 && <TrendingScene />}
      {scene === 2 && <SearchScene />}
      {scene === 3 && <RandomScene />}
    </AbsoluteFill>
  );
};
