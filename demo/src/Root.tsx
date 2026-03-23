import { Composition } from "remotion";
import { TerminalDemo } from "./TerminalDemo";

export const RemotionRoot: React.FC = () => {
  return (
    <Composition
      id="GithubDiscoverDemo"
      component={TerminalDemo}
      durationInFrames={360}
      fps={30}
      width={960}
      height={540}
    />
  );
};
