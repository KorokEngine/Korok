# 关于 Korok

Korok 取自塞尔达传说里面的森林精灵，在我的最初计划里面这应该是一个面向组件的简单易扩展的2D游戏引擎，引擎的设计启发自：bitsquid.blogspot.com 。
在一般的游戏引擎中都会有游戏对象的概念，比如：Unity 的 GameObject，Cocos2D 的 Sprite，在 korok 里面，是没有这个概念的，它的本质是一个ECS系统，GameObject
是通过一个 id索引各种相关组件组成的，现在设计了几种组件：

1. RenderComp 渲染一个可渲染对象
2. Transform  负责坐标和场景图
3. Animator   动画控制器
3. SourceComp 负责音频

--七夕的第二天

操千曲而后晓声，观千剑而后识器。在之前的设计中，渲染引擎的基础架构设计的不够好，导致后面的系统集成变得艰难。在阅读了大量的开源游戏引擎的渲染系统代码之后，对基础图形API接口的设计
有了新的想法。这样的API符合的特征是：

1. 无状态渲染提交(stateless)
2. 基于排序渲染(sort-based)
3. 易于实现多线程(multi-thread)
4. 兼容新的图形API

这部分的底层API，实现在 `gfx/bk` 叫做 'bk-api'。目前版本的实现搭建了基本的API框架，排序，多线程等暂时都没有实现。但是架构在此，以后添加上这些并非难事。
从现在的观点看来，这种API设计是比较主流的底层图形API设计方式，Ogre3D/Paradox/BitSquid 等大型3D开源引擎都采用了这种设计，此处的设计借鉴了 bgfx 的实现，可以说
是一个简版的 bgfx，但是移除了 bgfx 里面大量的资源管理功能。


--国庆的第二天

基于 bk-API，分别实现了一个用于渲染简单网格的 MeshRender 和 一个用于批量精灵纹理渲染的 BatchRender。Korok 的 Batch 系统，底层维护了8个VBO，最多可以生成128个Batch分组。
Batch接口采用了一个半自动的方案：提供一个 batch-id 字段，如果用户把一组图元标记为相同的 batch-id，那么他们可能会被batch在一起。还有另外两种常见方案：

1. 使用 SpriteBatch 提供手动的 Batch 接口
2. 引擎检测材质，自动根据不同材质进行合并

第一种略显过时，第二种常常会造成困扰，比如我们认为应该batch的场景却没有batch。而korok中采用半自动的方案，可以在无法batch时做出相应的提示，并且尽量按照用户希望
的方式batch，或许能提供更好的开发体验。

--2017/10/06

音频系统暂时可以工作了，这部分api实现在 `audio/ap` 叫做 'ap-API'。这是一套非常底层的API具体负责：

1. 统一的资源管理
2. 播放优先级决策

音频格式(编解码)的支持在 `audio/xx` 目录下实现，这里有两个已经支持的音频格式:wav 和 ogg. wav 一般用来测试和播放占用内存较少的音效，ogg 用来播放背景音乐之类。
此处并不打算支持 mp3 格式，相较来说无论文件大小还是保真程度 ogg 都是优于 mp3 的，所以应该尽量用 ogg 格式来编码游戏音频。编解码模块和底层音频播放系统是完全分离的，
我们通过上层 API 来配置音频的编解码，相互之间通过接口依赖，这样方便以后添加更多音频格式支持。

音频管理（内存管理）和播放部分还有一些功能没有实现(TODO)。

--2017/10/14

在 `gfx/dbg/` 目录下实现了 debug 绘制功能，这这种功能通常用来显示一些debug信息，比如fps,GPU和GPU状态监控等。目前可以用来在屏幕上打印矩形和字符串，这是一个自包含的
模块，可以说是一个微型的渲染系统，它不依赖GUI，也不依赖于引擎的渲染系统（直接基于bk-API）实现。我们的API使用了 *ImmediateMode GUI* 设计，这种 API 接口简单，却能够
实现强大的功能，对于打印一些临时信息来说再好不过，下面是一个段示例：

```
func OnLoop() {
	// draw something...
	dbg.Move(50, 50)
	dbg.Color(0xFFFF0000)
	dbg.DrawRect(0, 0, 50, 50)

	dbg.Move(300, 50)
	dbg.DrawStr("fps:60000!!")

	// advance frame
	dbg.NextFrame()
}
```
-- 2017/10/17

Korok的最终架构这几天确定了 - Comp/Table/System，我们的所有数据采用类似数据库表的概念来管理，所以对一个组件的增删改查类似于表的CRUD操作。我们的底层引擎也会使用同样的接口来操作数据
比如SpriteComp/MeshComp的各自存储在自己的表中。渲染引擎在渲染的时候也是读取当前的表查出所有可渲染对象然后再做绘制。此次没有功能上的更新，大多是运行时系统的重构，但是却
很重要，因为这次我们确定了整体架构和设计哲学，这会直接影响到以后的系统设计/开发。

-- 万圣节的第二天 🎃

可见性系统的设计还是很不完善，而且也没有太好的想法，暂时去掉这个功能，这样可以继续完善渲染系统。目前渲染系统的把Feature和Render分开了，Render是可以复用的基于
Shader的渲染工具，Feature是具体的渲染类型，比如SpriteComp/TextComp.. 它们都可以用BatchRender来渲染，只是使用上有细微的差别。由于 Golang 不支持泛型所以
只能采取间接层的方式--既Feature，来实现差异部分的功能。

-- 2017/12/04

当前已经可以使用 MeshComp 来渲染网格，用 SpriteComp 来渲染精灵图片，从接口API道底层渲染系统联调通过。中间写了一些临时代码，以后会去掉。批量渲染精灵的API如下:

```
    id, _ := assets.Texture.GetTexture("src/main/assets/ball.png")

	for i := 0; i < 50; i++ {
		face := korok.Entity.New()
		korok.Sprite.NewComp(face, assets.AsSubTexture(id))

		faceXF := korok.Transform.NewComp(face)

		x := float32(rand.Intn(480))
		y := float32(rand.Intn(320))
		faceXF.Position = mgl32.Vec2{x, y}
	}
```
今天是个值得庆祝的日子！😊
-- 2017/12/07

重新组织了包名和域名，现在可以使用 `go get korok.io/korok` 来下载引擎了(还没有启用HTTPS，需要加`-insecure`标记).
-- 2017/12/12
