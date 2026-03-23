import React from "react";
import {
  AbsoluteFill,
  useCurrentFrame,
  interpolate,
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

// Terminal window chrome
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
      <div
        style={{
          width: 12,
          height: 12,
          borderRadius: "50%",
          background: colors.red,
        }}
      />
      <div
        style={{
          width: 12,
          height: 12,
          borderRadius: "50%",
          background: colors.yellow,
        }}
      />
      <div
        style={{
          width: 12,
          height: 12,
          borderRadius: "50%",
          background: colors.green,
        }}
      />
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

// Typing animation
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
  const showCursor =
    charsToShow < text.length && (Math.floor(frame / 8) % 2 === 0 || true);

  return (
    <span style={{ color }}>
      {visible}
      {showCursor && charsToShow < text.length && (
        <span
          style={{
            background: colors.green,
            color: colors.bg,
            animation: "none",
          }}
        >
          {" "}
        </span>
      )}
    </span>
  );
};

// Fade-in line
const FadeInLine: React.FC<{
  children: React.ReactNode;
  startFrame: number;
}> = ({ children, startFrame }) => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();
  const opacity = spring({ frame: frame - startFrame, fps, config: { damping: 20 } });
  if (frame < startFrame) return null;
  return <div style={{ opacity }}>{children}</div>;
};

// Repo row component
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
    <div style={{ opacity: progress, transform: `translateX(${(1 - progress) * 20}px)` }}>
      <div style={{ display: "flex", gap: 8, alignItems: "baseline" }}>
        <span style={{ color: colors.gray, width: 30, textAlign: "right" }}>
          {rank}.
        </span>
        <span style={{ color: colors.green, width: 340 }}>{name}</span>
        <span style={{ color: colors.yellow }}>★ {stars}</span>
        {lang && (
          <span style={{ color: colors.cyan }}>[{lang}]</span>
        )}
      </div>
      <div style={{ color: colors.gray, paddingLeft: 38, fontSize: 13 }}>
        {desc}
      </div>
    </div>
  );
};

export const TerminalDemo: React.FC = () => {
  const frame = useCurrentFrame();

  // Scene 1: trending command (frames 0-180)
  // Scene 2: random command (frames 180-360)

  const scene = frame < 180 ? 1 : 2;

  return (
    <AbsoluteFill
      style={{
        background: "linear-gradient(135deg, #11111b 0%, #1e1e2e 100%)",
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      {scene === 1 && (
        <TerminalChrome>
          {/* Prompt */}
          <div>
            <span style={{ color: colors.green }}>$ </span>
            <TypedText
              text="github-discover trending -n 5"
              startFrame={10}
              color={colors.text}
            />
          </div>

          {/* Output */}
          <FadeInLine startFrame={40}>
            <div style={{ marginTop: 16 }}>
              <span style={{ color: colors.purple, fontWeight: "bold" }}>
                Trending Repositories (weekly)
              </span>
            </div>
            <div style={{ color: colors.border }}>
              {"─".repeat(70)}
            </div>
          </FadeInLine>

          <RepoRow
            rank={1}
            name="freeCodeCamp/freeCodeCamp"
            stars="438.7k"
            lang="TypeScript"
            desc="Open-source codebase and curriculum"
            startFrame={50}
          />
          <RepoRow
            rank={2}
            name="public-apis/public-apis"
            stars="414.4k"
            lang="Python"
            desc="A collective list of free APIs"
            startFrame={58}
          />
          <RepoRow
            rank={3}
            name="EbookFoundation/free-programming-books"
            stars="384.4k"
            lang="Python"
            desc="Freely available programming books"
            startFrame={66}
          />
          <RepoRow
            rank={4}
            name="kamranahmedse/developer-roadmap"
            stars="351.5k"
            lang="TypeScript"
            desc="Interactive roadmaps and guides for developers"
            startFrame={74}
          />
          <RepoRow
            rank={5}
            name="donnemartin/system-design-primer"
            stars="339.9k"
            lang="Python"
            desc="Learn how to design large-scale systems"
            startFrame={82}
          />

          <FadeInLine startFrame={92}>
            <div style={{ color: colors.border, marginTop: 4 }}>
              {"─".repeat(70)}
            </div>
            <div style={{ color: colors.gray, fontStyle: "italic", fontSize: 13 }}>
              Use --language to filter by language, --since to change time range
            </div>
          </FadeInLine>
        </TerminalChrome>
      )}

      {scene === 2 && (
        <TerminalChrome>
          {/* Prompt */}
          <div>
            <span style={{ color: colors.green }}>$ </span>
            <TypedText
              text="github-discover random"
              startFrame={185}
              color={colors.text}
            />
          </div>

          {/* Box output */}
          <FadeInLine startFrame={210}>
            <div
              style={{
                marginTop: 16,
                border: `1px solid ${colors.blue}`,
                borderRadius: 8,
                padding: "16px 20px",
              }}
            >
              <div
                style={{
                  color: colors.green,
                  fontWeight: "bold",
                  fontSize: 16,
                }}
              >
                torvalds/linux
              </div>
              <div
                style={{
                  color: colors.border,
                  margin: "8px 0",
                }}
              >
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

          <FadeInLine startFrame={240}>
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
      )}
    </AbsoluteFill>
  );
};
